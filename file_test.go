package backup

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FileEqual_Happy(t *testing.T) {
	f1 := file{
		name: "file1",
		size: 100,
	}

	f2 := file{
		name: "file1",
		size: 100,
	}

	assert.True(t, isEqual(f1, f2))
}

func Test_File_NotEqual_Name(t *testing.T) {
	f1 := file{
		name: "file1",
		size: 100,
	}

	f2 := file{
		name: "file2",
		size: 100,
	}

	assert.False(t, isEqual(f1, f2))
}

func Test_File_NotEqual_Size(t *testing.T) {
	f1 := file{
		name: "file1",
		size: 100,
	}

	f2 := file{
		name: "file1",
		size: 0,
	}

	assert.False(t, isEqual(f1, f2))
}

func Test_File_newFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "convertFromFileInfoDir")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)

	tmpFile, err := ioutil.TempFile(dir, "convertFromFileInfoFile")
	if err != nil {
		t.Fatal(err)
	}

	tempFileInfo, err := tmpFile.Stat()
	if err != nil {
		t.Fatal(err)
	}

	expected := file{
		name: tmpFile.Name(),
		size: tempFileInfo.Size(),
	}

	assert.Equal(t, expected, newFile(tmpFile.Name(), tempFileInfo.Size()))
}
