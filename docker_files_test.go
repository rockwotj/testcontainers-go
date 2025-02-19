package testcontainers_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCopyFileToContainer(t *testing.T) {
	ctx, cnl := context.WithTimeout(context.Background(), 30*time.Second)
	defer cnl()

	// copyFileOnCreate {
	absPath, err := filepath.Abs(filepath.Join(".", "testdata", "hello.sh"))
	if err != nil {
		t.Fatal(err)
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "docker.io/bash",
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/hello.sh",
					FileMode:          0o700,
				},
			},
			Cmd:        []string{"bash", "/hello.sh"},
			WaitingFor: wait.ForLog("done"),
		},
		Started: true,
	})
	// }

	require.NoError(t, err)
	require.NoError(t, container.Terminate(ctx))
}

func TestCopyFileToRunningContainer(t *testing.T) {
	ctx, cnl := context.WithTimeout(context.Background(), 30*time.Second)
	defer cnl()

	// Not using the assertations here to avoid leaking the library into the example
	// copyFileAfterCreate {
	waitForPath, err := filepath.Abs(filepath.Join(".", "testdata", "waitForHello.sh"))
	if err != nil {
		t.Fatal(err)
	}
	helloPath, err := filepath.Abs(filepath.Join(".", "testdata", "hello.sh"))
	if err != nil {
		t.Fatal(err)
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "docker.io/bash:5.2.26",
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      waitForPath,
					ContainerFilePath: "/waitForHello.sh",
					FileMode:          0o700,
				},
			},
			Cmd: []string{"bash", "/waitForHello.sh"},
		},
		Started: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = container.CopyFileToContainer(ctx, helloPath, "/scripts/hello.sh", 0o700)
	// }

	require.NoError(t, err)

	// Give some time to the wait script to catch the hello script being created
	err = wait.ForLog("done").WithStartupTimeout(2*time.Second).WaitUntilReady(ctx, container)
	require.NoError(t, err)

	require.NoError(t, container.Terminate(ctx))
}

func TestCopyDirectoryToContainer(t *testing.T) {
	ctx, cnl := context.WithTimeout(context.Background(), 30*time.Second)
	defer cnl()

	// Not using the assertations here to avoid leaking the library into the example
	// copyDirectoryToContainer {
	dataDirectory, err := filepath.Abs(filepath.Join(".", "testdata"))
	if err != nil {
		t.Fatal(err)
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "docker.io/bash",
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath: dataDirectory,
					// ContainerFile cannot create the parent directory, so we copy the scripts
					// to the root of the container instead. Make sure to create the container directory
					// before you copy a host directory on create.
					ContainerFilePath: "/",
					FileMode:          0o700,
				},
			},
			Cmd:        []string{"bash", "/testdata/hello.sh"},
			WaitingFor: wait.ForLog("done"),
		},
		Started: true,
	})
	// }

	require.NoError(t, err)
	require.NoError(t, container.Terminate(ctx))
}

func TestCopyDirectoryToRunningContainerAsFile(t *testing.T) {
	ctx, cnl := context.WithTimeout(context.Background(), 30*time.Second)
	defer cnl()

	// copyDirectoryToRunningContainerAsFile {
	dataDirectory, err := filepath.Abs(filepath.Join(".", "testdata"))
	if err != nil {
		t.Fatal(err)
	}
	waitForPath, err := filepath.Abs(filepath.Join(dataDirectory, "waitForHello.sh"))
	if err != nil {
		t.Fatal(err)
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "docker.io/bash",
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      waitForPath,
					ContainerFilePath: "/waitForHello.sh",
					FileMode:          0o700,
				},
			},
			Cmd: []string{"bash", "/waitForHello.sh"},
		},
		Started: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// as the container is started, we can create the directory first
	_, _, err = container.Exec(ctx, []string{"mkdir", "-p", "/scripts"})
	if err != nil {
		t.Fatal(err)
	}

	// because the container path is a directory, it will use the copy dir method as fallback
	err = container.CopyFileToContainer(ctx, dataDirectory, "/scripts", 0o700)
	if err != nil {
		t.Fatal(err)
	}
	// }

	require.NoError(t, err)
	require.NoError(t, container.Terminate(ctx))
}

func TestCopyDirectoryToRunningContainerAsDir(t *testing.T) {
	ctx, cnl := context.WithTimeout(context.Background(), 30*time.Second)
	defer cnl()

	// Not using the assertations here to avoid leaking the library into the example
	// copyDirectoryToRunningContainerAsDir {
	waitForPath, err := filepath.Abs(filepath.Join(".", "testdata", "waitForHello.sh"))
	if err != nil {
		t.Fatal(err)
	}
	dataDirectory, err := filepath.Abs(filepath.Join(".", "testdata"))
	if err != nil {
		t.Fatal(err)
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "docker.io/bash",
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      waitForPath,
					ContainerFilePath: "/waitForHello.sh",
					FileMode:          0o700,
				},
			},
			Cmd: []string{"bash", "/waitForHello.sh"},
		},
		Started: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// as the container is started, we can create the directory first
	_, _, err = container.Exec(ctx, []string{"mkdir", "-p", "/scripts"})
	if err != nil {
		t.Fatal(err)
	}

	err = container.CopyDirToContainer(ctx, dataDirectory, "/scripts", 0o700)
	if err != nil {
		t.Fatal(err)
	}
	// }

	require.NoError(t, err)
	require.NoError(t, container.Terminate(ctx))
}
