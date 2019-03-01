package hermetic_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/changes/hermetic"
)

var _ = Describe("copying files and directories", func() {
	Describe("RecursiveCopy()", func() {
		Context("when the source is file and the destination is a file", func() {
			It("copies the file, creating any nonexistant directories", func() {
				tmpDir, err := ioutil.TempDir("", "hermetic-recursive-copy-test")
				Expect(err).ToNot(HaveOccurred(), "expected no error creating a temporary directory")
				defer os.RemoveAll(tmpDir)

				filename := "test_file.md"
				err = ioutil.WriteFile(filepath.Join(tmpDir, filename), []byte("some test content"), os.ModePerm)
				Expect(err).ToNot(HaveOccurred(), "expected no error writing a test temporary file")

				newDirectory := "new_directory/"

				err = hermetic.RecursiveCopy(filepath.Join(tmpDir, filename), filepath.Join(tmpDir, newDirectory, filename))
				Expect(err).ToNot(HaveOccurred(), "expected no error when copying the test file content to a non existant directory")
				Expect(filepath.Join(tmpDir, newDirectory, filename)).To(BeARegularFile(), "expected to find copied file in previously non existant directory")
			})
		})

		Context("when the source is a file and the destination is a directory", func() {
			It("copies the file to the directory, creating any nonexistant directories", func() {
				tmpDir, err := ioutil.TempDir("", "hermetic-recursive-copy-test")
				Expect(err).ToNot(HaveOccurred(), "expected no error creating a temporary directory")
				defer os.RemoveAll(tmpDir)

				filename := "test_file.md"
				err = ioutil.WriteFile(filepath.Join(tmpDir, filename), []byte("some test content"), os.ModePerm)
				Expect(err).ToNot(HaveOccurred(), "expected no error writing a test temporary file")

				newDirectory := "new_directory/"

				err = hermetic.RecursiveCopy(filepath.Join(tmpDir, filename), filepath.Join(tmpDir, newDirectory)+string(filepath.Separator))
				Expect(err).ToNot(HaveOccurred(), "expected no error when copying the test file content to a non existant directory")
				Expect(filepath.Join(tmpDir, newDirectory, filename)).To(BeARegularFile(), "expected to find copied file in previously non existant directory")
			})
		})

		Context("when the source is a directory and the destination is a directory", func() {
			It("recursively copies files and directories, creating those which don't exist", func() {
				tmpDir, err := ioutil.TempDir("", "hermetic-recursive-copy-test")
				Expect(err).ToNot(HaveOccurred(), "expected no error creating a temporary directory")
				defer os.RemoveAll(tmpDir)

				filename := "test_file.md"
				testContent := []byte("test content\n")

				newTreeOne := "a/b/c/"
				newTreeTwo := "a/b/c/d/"
				newTreeThree := "a/b/c/d/e/f/"

				targetLocation := "new_directory/"

				err = os.MkdirAll(filepath.Join(tmpDir, newTreeOne), os.ModePerm)
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating test directories")

				err = os.MkdirAll(filepath.Join(tmpDir, newTreeTwo), os.ModePerm)
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating test directories")

				err = os.MkdirAll(filepath.Join(tmpDir, newTreeThree), os.ModePerm)
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating test directories")

				err = ioutil.WriteFile(filepath.Join(tmpDir, newTreeOne, filename), testContent, os.ModePerm)
				Expect(err).ToNot(HaveOccurred(), "expected no error when writing a test file")

				err = ioutil.WriteFile(filepath.Join(tmpDir, newTreeTwo, filename), testContent, os.ModePerm)
				Expect(err).ToNot(HaveOccurred(), "expected no error when writing a test file")

				err = ioutil.WriteFile(filepath.Join(tmpDir, newTreeThree, filename), testContent, os.ModePerm)
				Expect(err).ToNot(HaveOccurred(), "expected no error when writing a test file")

				err = hermetic.RecursiveCopy(filepath.Join(tmpDir, "a"), filepath.Join(tmpDir, targetLocation)+string(filepath.Separator))
				Expect(err).ToNot(HaveOccurred(), "expected no error when copying the test file content to a non existant directory")

				Expect(filepath.Join(tmpDir, targetLocation, newTreeOne, filename)).To(BeARegularFile(), "expected to find copied file in previously non existant directory")

				copiedBytes, err := ioutil.ReadFile(filepath.Join(tmpDir, targetLocation, newTreeOne, filename))
				Expect(err).ToNot(HaveOccurred(), "expected no error when opening a copied file which exists")
				Expect(string(copiedBytes)).To(Equal(string(testContent)), "expected copied content to equal test content")

				Expect(filepath.Join(tmpDir, targetLocation, newTreeTwo, filename)).To(BeARegularFile(), "expected to find copied file in previously non existant directory")
				copiedBytes, err = ioutil.ReadFile(filepath.Join(tmpDir, targetLocation, newTreeTwo, filename))
				Expect(err).ToNot(HaveOccurred(), "expected no error when opening a copied file which exists")
				Expect(string(copiedBytes)).To(Equal(string(testContent)), "expected copied content to equal test content")

				Expect(filepath.Join(tmpDir, targetLocation, newTreeThree, filename)).To(BeARegularFile(), "expected to find copied file in previously non existant directory")
				copiedBytes, err = ioutil.ReadFile(filepath.Join(tmpDir, targetLocation, newTreeThree, filename))
				Expect(err).ToNot(HaveOccurred(), "expected no error when opening a copied file which exists")
				Expect(string(copiedBytes)).To(Equal(string(testContent)), "expected copied content to equal test content")
			})
		})

		Context("when the source is a directory and the destination is a file", func() {
			It("returns an error", func() {
				tmpDir, err := ioutil.TempDir("", "hermetic-recursive-copy-test")
				Expect(err).ToNot(HaveOccurred(), "expected no error creating a temporary directory")
				defer os.RemoveAll(tmpDir)

				target := "cant.iso"
				err = hermetic.RecursiveCopy(tmpDir, target)
				Expect(err).To(MatchError(fmt.Sprintf("cannot copy directory: %s into file: %s", tmpDir, target)), "expected error to contain failure to copy directory into file")
			})
		})

	})
})
