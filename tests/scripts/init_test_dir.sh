#!/bin/bash

TEST_DIR="test_dir"
rm -rf $TEST_DIR
mkdir -p $TEST_DIR/dir1
mkdir -p $TEST_DIR/dir2/subdir

echo "content1" > $TEST_DIR/dir1/file1.txt
echo "content2" > $TEST_DIR/dir1/file2.log
echo "content3" > $TEST_DIR/dir2/file3.txt
echo "content4" > $TEST_DIR/dir2/subdir/file4.md

ln -s dir1/file1.txt $TEST_DIR/symlink1
ln -s nonexistent.txt $TEST_DIR/broken_symlink
