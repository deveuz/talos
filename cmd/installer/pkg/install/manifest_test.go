// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package install_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/talos-systems/go-blockdevice/blockdevice"
	"github.com/talos-systems/go-blockdevice/blockdevice/loopback"
	"github.com/talos-systems/go-blockdevice/blockdevice/partition/gpt"
	"github.com/talos-systems/go-blockdevice/blockdevice/util"

	"github.com/talos-systems/talos/cmd/installer/pkg/install"
	"github.com/talos-systems/talos/internal/app/machined/pkg/runtime"
	"github.com/talos-systems/talos/internal/pkg/mount"
	"github.com/talos-systems/talos/pkg/machinery/constants"
	"github.com/talos-systems/talos/pkg/makefs"
)

// Some tests in this package cannot be run under buildkit, as buildkit doesn't propagate partition devices
// like /dev/loopXpY into the sandbox. To run the tests on your local computer, do the following:
//
//	 sudo go test -v --count 1 ./cmd/installer/pkg/install/

type manifestSuite struct {
	suite.Suite

	disk           *os.File
	loopbackDevice *os.File
}

const (
	diskSize    = 4 * 1024 * 1024 * 1024 // 4 GiB
	lbaSize     = 512
	gptReserved = 67
)

const (
	legacyBootSize      = 512 * install.MiB
	legacyEphemeralSize = diskSize - legacyBootSize - gptReserved*lbaSize
)

func TestManifestSuite(t *testing.T) {
	suite.Run(t, new(manifestSuite))
}

func (suite *manifestSuite) SetupTest() {
	suite.skipIfNotRoot()

	var err error

	suite.disk, err = ioutil.TempFile("", "talos")
	suite.Require().NoError(err)

	suite.Require().NoError(suite.disk.Truncate(diskSize))

	suite.loopbackDevice, err = loopback.NextLoopDevice()
	suite.Require().NoError(err)

	suite.T().Logf("Using %s", suite.loopbackDevice.Name())

	suite.Require().NoError(loopback.Loop(suite.loopbackDevice, suite.disk))

	suite.Require().NoError(loopback.LoopSetReadWrite(suite.loopbackDevice))
}

func (suite *manifestSuite) TearDownTest() {
	if suite.loopbackDevice != nil {
		suite.Assert().NoError(loopback.Unloop(suite.loopbackDevice))
	}

	if suite.disk != nil {
		suite.Assert().NoError(os.Remove(suite.disk.Name()))
		suite.Assert().NoError(suite.disk.Close())
	}
}

func (suite *manifestSuite) skipUnderBuildkit() {
	hostname, _ := os.Hostname() //nolint: errcheck

	if hostname == "buildkitsandbox" {
		suite.T().Skip("test not supported under buildkit as partition devices are not propagated from /dev")
	}
}

func (suite *manifestSuite) skipIfNotRoot() {
	if os.Getuid() != 0 {
		suite.T().Skip("can't run the test as non-root")
	}
}

func (suite *manifestSuite) verifyBlockdevice(manifest *install.Manifest, current, next string, verifyConfigPersistence, verifyEphemeralPersistence, upgradeFromLegacy bool) {
	bd, err := blockdevice.Open(suite.loopbackDevice.Name())
	suite.Require().NoError(err)

	defer bd.Close() //nolint: errcheck

	table, err := bd.PartitionTable()
	suite.Require().NoError(err)

	// verify partition table

	suite.Assert().Len(table.Partitions().Items(), 6)

	part := table.Partitions().Items()[0]
	suite.Assert().Equal(install.EFISystemPartition, strings.ToUpper(part.Type.String()))
	suite.Assert().Equal(constants.EFIPartitionLabel, part.Name)
	suite.Assert().EqualValues(0, part.Attributes)
	suite.Assert().EqualValues(install.EFISize/lbaSize, part.Length())

	part = table.Partitions().Items()[1]
	suite.Assert().Equal(install.BIOSBootPartition, strings.ToUpper(part.Type.String()))
	suite.Assert().Equal(constants.BIOSGrubPartitionLabel, part.Name)
	suite.Assert().EqualValues(4, part.Attributes)
	suite.Assert().EqualValues(install.BIOSGrubSize/lbaSize, part.Length())

	part = table.Partitions().Items()[2]
	suite.Assert().Equal(install.LinuxFilesystemData, strings.ToUpper(part.Type.String()))
	suite.Assert().Equal(constants.BootPartitionLabel, part.Name)
	suite.Assert().EqualValues(0, part.Attributes)
	suite.Assert().EqualValues(install.BootSize/lbaSize, part.Length())

	part = table.Partitions().Items()[3]
	suite.Assert().Equal(install.LinuxFilesystemData, strings.ToUpper(part.Type.String()))
	suite.Assert().Equal(constants.MetaPartitionLabel, part.Name)
	suite.Assert().EqualValues(0, part.Attributes)
	suite.Assert().EqualValues(install.MetaSize/lbaSize, part.Length())

	part = table.Partitions().Items()[4]
	suite.Assert().Equal(install.LinuxFilesystemData, strings.ToUpper(part.Type.String()))
	suite.Assert().Equal(constants.StatePartitionLabel, part.Name)
	suite.Assert().EqualValues(0, part.Attributes)

	if !upgradeFromLegacy {
		suite.Assert().EqualValues(install.StateSize/lbaSize, part.Length())
	} else {
		suite.Assert().EqualValues((diskSize-legacyEphemeralSize-install.EFISize-install.BIOSGrubSize-install.BootSize-install.MetaSize)/lbaSize-gptReserved, part.Length())
	}

	part = table.Partitions().Items()[5]
	suite.Assert().Equal(install.LinuxFilesystemData, strings.ToUpper(part.Type.String()))
	suite.Assert().Equal(constants.EphemeralPartitionLabel, part.Name)
	suite.Assert().EqualValues(0, part.Attributes)

	if !upgradeFromLegacy {
		suite.Assert().EqualValues((diskSize-install.EFISize-install.BIOSGrubSize-install.BootSize-install.MetaSize-install.StateSize)/lbaSize-gptReserved, part.Length())
	} else {
		suite.Assert().EqualValues(legacyEphemeralSize/lbaSize, part.Length())
	}

	suite.Assert().NoError(bd.Close())

	// query mount points directly for the device

	mountpoints, err := mount.SystemMountPointsForDevice(suite.loopbackDevice.Name())
	suite.Require().NoError(err)

	suite.Assert().Equal(4, mountpoints.Len())

	// verify filesystems by mounting and unmounting

	tempDir, err := ioutil.TempDir("", "talos")
	suite.Require().NoError(err)

	defer func() {
		suite.Assert().NoError(os.RemoveAll(tempDir))
	}()

	mountpoints, err = manifest.SystemMountpoints()
	suite.Require().NoError(err)

	suite.Assert().Equal(4, mountpoints.Len())

	suite.Require().NoError(mount.PrefixMountTargets(mountpoints, tempDir))

	err = mount.Mount(mountpoints)
	suite.Require().NoError(err)

	defer func() {
		suite.Assert().NoError(mount.Unmount(mountpoints))
	}()

	metaPath := fmt.Sprintf("%sp%d", suite.loopbackDevice.Name(), table.Partitions().Items()[3].Number)

	if verifyConfigPersistence {
		suite.Assert().FileExists(filepath.Join(tempDir, "system", "state", "config.yaml"))
	}

	if verifyEphemeralPersistence {
		suite.Assert().FileExists(filepath.Join(tempDir, "var", "content"))
	}

	if current != "" {
		// verify that current was preserved
		suite.Assert().DirExists(filepath.Join(tempDir, "boot", current))

		suite.Assert().FileExists(filepath.Join(tempDir, "boot", current, "kernel"))

		buf := make([]byte, len(current))

		f, err := os.Open(metaPath)
		suite.Require().NoError(err)

		_, err = io.ReadFull(f, buf)
		suite.Require().NoError(err)

		suite.Assert().Equal(current, string(buf))

		suite.Assert().NoError(f.Close())
	}

	if next != "" {
		suite.Assert().NoError(os.MkdirAll(filepath.Join(tempDir, "boot", next), 0o700))
		suite.Assert().NoError(ioutil.WriteFile(filepath.Join(tempDir, "boot", next, "kernel"), []byte("LINUX!"), 0o660))
		suite.Assert().NoError(ioutil.WriteFile(filepath.Join(tempDir, "system", "state", "config.yaml"), []byte("#!yaml"), 0o660))

		buf := []byte(next)

		f, err := os.OpenFile(metaPath, os.O_WRONLY, 0)
		suite.Require().NoError(err)

		_, err = f.Write(buf)
		suite.Require().NoError(err)

		suite.Assert().NoError(f.Close())
	}

	suite.Assert().NoError(ioutil.WriteFile(filepath.Join(tempDir, "var", "content"), []byte("data"), 0o600))
}

func (suite *manifestSuite) TestExecuteManifestClean() {
	suite.skipUnderBuildkit()

	manifest, err := install.NewManifest("A", runtime.SequenceInstall, false, &install.Options{
		Disk:       suite.loopbackDevice.Name(),
		Bootloader: true,
		Force:      true,
		Board:      constants.BoardNone,
	})
	suite.Require().NoError(err)

	suite.Assert().NoError(manifest.Execute())

	suite.verifyBlockdevice(manifest, "", "A", false, false, false)
}

func (suite *manifestSuite) TestExecuteManifestForce() {
	suite.skipUnderBuildkit()

	manifest, err := install.NewManifest("A", runtime.SequenceInstall, false, &install.Options{
		Disk:       suite.loopbackDevice.Name(),
		Bootloader: true,
		Force:      true,
		Board:      constants.BoardNone,
	})
	suite.Require().NoError(err)

	suite.Assert().NoError(manifest.Execute())

	suite.verifyBlockdevice(manifest, "", "A", false, false, false)

	// reinstall

	manifest, err = install.NewManifest("B", runtime.SequenceUpgrade, true, &install.Options{
		Disk:       suite.loopbackDevice.Name(),
		Bootloader: true,
		Force:      true,
		Zero:       true,
		Board:      constants.BoardNone,
	})
	suite.Require().NoError(err)

	suite.Assert().NoError(manifest.Execute())

	suite.verifyBlockdevice(manifest, "A", "B", true, false, false)
}

func (suite *manifestSuite) TestExecuteManifestPreserve() {
	suite.skipUnderBuildkit()

	manifest, err := install.NewManifest("A", runtime.SequenceInstall, false, &install.Options{
		Disk:       suite.loopbackDevice.Name(),
		Bootloader: true,
		Force:      true,
		Board:      constants.BoardNone,
	})
	suite.Require().NoError(err)

	suite.Assert().NoError(manifest.Execute())

	suite.verifyBlockdevice(manifest, "", "A", false, false, false)

	// reinstall

	manifest, err = install.NewManifest("B", runtime.SequenceUpgrade, true, &install.Options{
		Disk:       suite.loopbackDevice.Name(),
		Bootloader: true,
		Force:      false,
		Board:      constants.BoardNone,
	})
	suite.Require().NoError(err)

	suite.Assert().NoError(manifest.Execute())

	suite.verifyBlockdevice(manifest, "A", "B", true, true, false)
}

func (suite *manifestSuite) TestExecuteManifestLegacyForce() {
	suite.skipUnderBuildkit()

	suite.createTalosLegacyLayout()

	// upgrade with force

	manifest, err := install.NewManifest("A", runtime.SequenceUpgrade, true, &install.Options{
		Disk:       suite.loopbackDevice.Name(),
		Bootloader: true,
		Force:      true,
		Board:      constants.BoardNone,
	})
	suite.Require().NoError(err)

	suite.Assert().NoError(manifest.Execute())

	suite.verifyBlockdevice(manifest, "", "", true, false, false)
}

func (suite *manifestSuite) TestExecuteManifestLegacyPreserve() {
	suite.skipUnderBuildkit()

	suite.createTalosLegacyLayout()

	// upgrade with preserve

	manifest, err := install.NewManifest("A", runtime.SequenceUpgrade, true, &install.Options{
		Disk:       suite.loopbackDevice.Name(),
		Bootloader: true,
		Force:      false,
		Board:      constants.BoardNone,
	})
	suite.Require().NoError(err)

	suite.Assert().NoError(manifest.Execute())

	suite.verifyBlockdevice(manifest, "", "", true, true, true)
}

func (suite *manifestSuite) TestTargetInstall() {
	// Create Temp dirname for mountpoint
	dir, err := ioutil.TempDir("", "talostest")
	suite.Require().NoError(err)

	// nolint: errcheck
	defer os.RemoveAll(dir)

	// Create a tempfile for local copy
	src, err := ioutil.TempFile(dir, "example")
	suite.Require().NoError(err)

	suite.Require().NoError(src.Close())

	dst := filepath.Join(dir, "dest")

	// Attempt to download and copy files
	target := &install.Target{
		Assets: []*install.Asset{
			{
				Source:      src.Name(),
				Destination: dst,
			},
		},
	}

	suite.Require().NoError(target.Save())

	for _, expectedFile := range target.Assets {
		// Verify copied file is at the appropriate location.
		_, err := os.Stat(expectedFile.Destination)
		suite.Require().NoError(err)
	}
}

func (suite *manifestSuite) createTalosLegacyLayout() {
	bd, err := blockdevice.Open(suite.loopbackDevice.Name())
	suite.Require().NoError(err)

	defer bd.Close() //nolint: errcheck

	// create Talos 0.6 partitions
	table, err := gpt.New(bd.Device())
	suite.Require().NoError(err)

	partBoot, err := table.Add(512*install.MiB,
		gpt.WithLegacyBIOSBootableAttribute(true),
		gpt.WithPartitionName(constants.LegacyBootPartitionLabel),
		gpt.WithPartitionType("28732AC1-1FF8-D211-BA4B-00A0C93EC93B"),
	)
	suite.Require().NoError(err)

	partEphemeral, err := table.Add(0,
		gpt.WithPartitionName(constants.EphemeralPartitionLabel),
		gpt.WithPartitionType("0FC63DAF-8483-4772-8E79-3D69D8477DE4"),
		gpt.WithMaximumSize(true),
	)
	suite.Require().NoError(err)

	suite.Require().NoError(table.Write())

	suite.Require().NoError(bd.Close())

	// format partitions
	partBootPath, err := util.PartPath(suite.loopbackDevice.Name(), int(partBoot.Number))
	suite.Require().NoError(err)

	suite.Require().NoError(makefs.VFAT(partBootPath))

	partEphemeralPath, err := util.PartPath(suite.loopbackDevice.Name(), int(partEphemeral.Number))
	suite.Require().NoError(err)

	suite.Require().NoError(makefs.XFS(partEphemeralPath, makefs.WithLabel(constants.EphemeralPartitionLabel)))

	// mount partitions temporarily and fill with data
	tempDir, err := ioutil.TempDir("", "talos")
	suite.Require().NoError(err)

	defer func() {
		suite.Assert().NoError(os.RemoveAll(tempDir))
	}()

	mountpoints := mount.NewMountPoints()
	mountpoints.Set(constants.LegacyBootPartitionLabel, mount.NewMountPoint(partBootPath, filepath.Join(tempDir, "boot"), install.FilesystemTypeVFAT, 0, ""))
	mountpoints.Set(constants.EphemeralPartitionLabel, mount.NewMountPoint(partEphemeralPath, filepath.Join(tempDir, "var"), install.FilesystemTypeXFS, 0, ""))

	err = mount.Mount(mountpoints)
	suite.Require().NoError(err)

	defer func() {
		suite.Assert().NoError(mount.Unmount(mountpoints))
	}()

	suite.Assert().NoError(ioutil.WriteFile(filepath.Join(tempDir, "boot", "config.yaml"), []byte("#!yaml"), 0o600))
	suite.Assert().NoError(ioutil.WriteFile(filepath.Join(tempDir, "var", "content"), []byte("data"), 0o600))
}
