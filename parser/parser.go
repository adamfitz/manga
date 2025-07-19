package parser

import (
	"fmt"
	"io"
	"io/fs"
	//"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func MapKeys(m map[string]any) []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// search the value of a nested map by a list of keys.  Send each key is as an individual string
func NestedMapValue(m map[string]any, keys ...string) (any, error) {
	var current any = m

	for _, key := range keys {
		// If current is a map, assert it to map[string]interface{}
		if mapValue, ok := current.(map[string]any); ok {
			current = mapValue[key]
		} else {
			return nil, fmt.Errorf("key '%s' not found or not a map at level", key)
		}
	}

	return current, nil
}

// Return a list of all directories from the provided rootDir.
// Optionally pass an exclusion list to skip certain directories, the slice is optional indicated by the variadic parameter.
func DirList(rootDir string, exclusionList ...string) ([]string, error) {

	// Convert exclusionList slice to a map for fast lookup
	exclusions := make(map[string]struct{}, len(exclusionList))
	for _, name := range exclusionList {
		exclusions[name] = struct{}{}
	}

	dirList := make([]string, 0)

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if _, skip := exclusions[entry.Name()]; !skip {
				dirList = append(dirList, entry.Name())
			}
		}
	}

	return dirList, nil
}

// CopyDir copies a directory recursively
func CopyDir(src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relativePath)

		if path == src {
			return nil
		}
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})

	return err
}

func DirsAreEqual(src, dst string) (bool, error) {
	srcEntries := map[string]fs.FileInfo{}
	dstEntries := map[string]fs.FileInfo{}

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		srcEntries[rel] = info
		return nil
	})
	if err != nil {
		return false, err
	}

	err = filepath.Walk(dst, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		rel, err := filepath.Rel(dst, path)
		if err != nil {
			return err
		}
		dstEntries[rel] = info
		return nil
	})
	if err != nil {
		return false, err
	}

	// Compare file count
	if len(srcEntries) != len(dstEntries) {
		return false, nil
	}

	// Compare file names and sizes
	for rel, srcInfo := range srcEntries {
		dstInfo, exists := dstEntries[rel]
		if !exists {
			return false, nil
		}
		if srcInfo.IsDir() != dstInfo.IsDir() {
			return false, nil
		}
		if !srcInfo.IsDir() && srcInfo.Size() != dstInfo.Size() {
			return false, nil
		}
	}

	return true, nil
}

// Return elements from sourceSlice that are not present in targetSlice.
func FindUniqueStrings(sourceSlice, targetSlice []string) []string {
	targetMap := make(map[string]bool)
	// normalise strings lowercase
	for _, s := range targetSlice {
		clean := strings.ToLower(strings.TrimSpace(s))
		targetMap[clean] = true
	}

	var uniqueStrings []string
	for _, s := range sourceSlice {
		clean := strings.ToLower(strings.TrimSpace(s))
		if !targetMap[clean] {
			uniqueStrings = append(uniqueStrings, s) // Return original string
		}
	}
	return uniqueStrings
}

// Merge multiple string slices, normalizes (lowercase + trim) and removes duplicates while returning the cleaned list.
func NormalizeAndDeduplicate(slices ...[]string) []string {
	nameMap := make(map[string]bool)

	for _, slice := range slices {
		for _, name := range slice {
			clean := strings.ToLower(strings.TrimSpace(name))
			nameMap[clean] = true
		}
	}

	var uniqueList []string
	for name := range nameMap {
		uniqueList = append(uniqueList, name)
	}
	return uniqueList
}

// Compares the directory names against mangadexList and mangaList (table name column), writes the manga name and the table its foudn in to a file (or both for both tables)
func WriteMissingDirsWithSourceTags(outputFile string, namesToWrite, mangadexList, mangaList []string) error {
	mangadexMap := make(map[string]bool)
	for _, name := range mangadexList {
		mangadexMap[strings.ToLower(strings.TrimSpace(name))] = true
	}

	mangaMap := make(map[string]bool)
	for _, name := range mangaList {
		mangaMap[strings.ToLower(strings.TrimSpace(name))] = true
	}

	// Sort the names alphabetically
	sort.Strings(namesToWrite)

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, name := range namesToWrite {
		clean := strings.ToLower(strings.TrimSpace(name))
		var tags []string

		if mangadexMap[clean] {
			tags = append(tags, "[MANGADEX]")
		}
		if mangaMap[clean] {
			tags = append(tags, "[MANGA]")
		}

		line := name + "\t" + strings.Join(tags, " ") + "\n"
		if _, err := f.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}

// Compare directory names to mangadex and manga table name column,and writes those missing from the DB to the output file with appropriate tags.
func WriteMissingTableEntriesWithSourceTags(outputFile string, dirNames, mangadexList, mangaList []string) error {
	mangadexMap := make(map[string]bool)
	for _, name := range mangadexList {
		mangadexMap[strings.ToLower(strings.TrimSpace(name))] = true
	}

	mangaMap := make(map[string]bool)
	for _, name := range mangaList {
		mangaMap[strings.ToLower(strings.TrimSpace(name))] = true
	}

	// Filter names missing from both tables
	var missingDirs []string
	for _, dir := range dirNames {
		clean := strings.ToLower(strings.TrimSpace(dir))
		if !mangadexMap[clean] && !mangaMap[clean] {
			missingDirs = append(missingDirs, dir)
		}
	}

	// Sort alphabetically
	sort.Strings(missingDirs)

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, dir := range missingDirs {
		if _, err := f.WriteString(dir + "\n"); err != nil {
			return err
		}
	}

	return nil
}
