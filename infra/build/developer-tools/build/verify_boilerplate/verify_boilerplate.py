#!/usr/bin/env python

# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# Verifies that all source files contain the necessary copyright boilerplate
# snippet.
# This is based on existing work
# https://github.com/kubernetes/test-infra/blob/master/hack
# /verify_boilerplate.py

# Please note that this file was generated from
# [terraform-google-module-template](https://github.com/terraform-google-modules/terraform-google-module-template).
# Please make sure to contribute relevant changes upstream!
from __future__ import print_function
import argparse
import glob
import os
import re
import sys


def get_args():
    """Parses command line arguments.

    Configures and runs argparse.ArgumentParser to extract command line
    arguments.

    Returns:
        An argparse.Namespace containing the arguments parsed from the
        command line
    """
    parser = argparse.ArgumentParser()
    parser.add_argument("filenames",
                        help="list of files to check, "
                             "all files if unspecified",
                        nargs='*')
    rootdir = os.path.abspath(os.getcwd())
    parser.add_argument(
        "--rootdir",
        default=rootdir,
        help="root directory to examine")

    default_boilerplate_dir = os.path.join(os.path.dirname(__file__),
                                           "boilerplate")
    parser.add_argument("--boilerplate-dir", default=default_boilerplate_dir)
    return parser.parse_args()


def get_refs(ARGS):
    """Converts the directory of boilerplate files into a map keyed by file
    extension.

    Reads each boilerplate file's contents into an array, then adds that array
    to a map keyed by the file extension.

    Returns:
        A map of boilerplate lines, keyed by file extension. For example,
        boilerplate.py.txt would result in the k,v pair {".py": py_lines} where
        py_lines is an array containing each line of the file.
    """
    refs = {}

    # Find and iterate over the absolute path for each boilerplate template
    for path in glob.glob(os.path.join(
            ARGS.boilerplate_dir,
            "boilerplate.*.txt")):
        extension = os.path.basename(path).split(".")[1]
        ref_file = open(path, 'r')
        ref = ref_file.read().splitlines()
        ref_file.close()
        refs[extension] = ref
    return refs


# pylint: disable=too-many-locals
def has_valid_header(filename, refs):
    """Test whether a file has the correct boilerplate header.

    Tests each file against the boilerplate stored in refs for that file type
    (based on extension), or by the entire filename (eg Dockerfile, Makefile).
    Some heuristics are applied to remove build tags and shebangs, but little
    variance in header formatting is tolerated.

    Args:
        filename: A string containing the name of the file to test
        refs: A map of boilerplate headers, keyed by file extension

    Returns:
        True if the file has the correct boilerplate header, otherwise returns
        False.
    """
    try:
        with open(filename, 'r') as fp:  # pylint: disable=invalid-name
            data = fp.read()
    except IOError:
        print(filename)
        return False
    basename = os.path.basename(filename)
    extension = get_file_extension(filename)
    if extension:
        ref = refs[extension]
    else:
        ref = refs[basename]
    data = data.splitlines()
    pattern_len = len(ref)
    # if our test file is smaller than the reference it surely fails!
    if pattern_len > len(data):
        return False
    copyright_regex = re.compile("Copyright 20\\d\\d")
    substitute_string = "Copyright YYYY"
    copyright_is_found = False
    j = 0
    for datum in data:
        # if it's a copyright line
        if not copyright_is_found and copyright_regex.search(datum):
            copyright_is_found = True
            # replace the actual year (e.g. 2019) with "YYYY" placeholder
            # used in a boilerplate
            datum = copyright_regex.sub(substitute_string, datum)
        if datum == ref[j]:
            j = j + 1
        else:
            j = 0
        if j == pattern_len:
            return copyright_is_found
    return copyright_is_found and j == pattern_len


def get_file_extension(filename):
    """Extracts the extension part of a filename.

    Identifies the extension as everything after the last period in filename.

    Args:
        filename: string containing the filename

    Returns:
        A string containing the extension in lowercase
    """
    return os.path.splitext(filename)[1].split(".")[-1].lower()


# These directories will be omitted from header checks
SKIPPED_DIRS = [
    'Godeps', 'third_party', '_gopath', '_output',
    '.git', 'vendor', '__init__.py', 'node_modules'
]


def normalize_files(files):
    """Extracts the files that require boilerplate checking from the files
    argument.

    A new list will be built. Each path from the original files argument will
    be added unless it is within one of SKIPPED_DIRS. All relative paths will
    be converted to absolute paths by prepending the root_dir path parsed from
    the command line, or its default value.

    Args:
        files: a list of file path strings

    Returns:
        A modified copy of the files list where any any path in a skipped
        directory is removed, and all paths have been made absolute.
    """
    newfiles = []
    for pathname in files:
        if any(x in pathname for x in SKIPPED_DIRS):
            continue
        newfiles.append(pathname)
    for idx, pathname in enumerate(newfiles):
        if not os.path.isabs(pathname):
            newfiles[idx] = os.path.join(ARGS.rootdir, pathname)
    return newfiles


def get_files(extensions, ARGS):
    """Generates a list of paths whose boilerplate should be verified.

    If a list of file names has been provided on the command line, it will be
    treated as the initial set to search. Otherwise, all paths within rootdir
    will be discovered and used as the initial set.

    Once the initial set of files is identified, it is normalized via
    normalize_files() and further stripped of any file name whose extension is
    not in extensions.

    Args:
        extensions: a list of file extensions indicating which file types
                    should have their boilerplate verified

    Returns:
        A list of absolute file paths
    """
    files = []
    if ARGS.filenames:
        files = ARGS.filenames
    else:
        for root, dirs, walkfiles in os.walk(ARGS.rootdir):
            # don't visit certain dirs. This is just a performance improvement
            # as we would prune these later in normalize_files(). But doing it
            # cuts down the amount of filesystem walking we do and cuts down
            # the size of the file list
            for dpath in SKIPPED_DIRS:
                if dpath in dirs:
                    dirs.remove(dpath)
            for name in walkfiles:
                pathname = os.path.join(root, name)
                files.append(pathname)
    files = normalize_files(files)
    outfiles = []
    for pathname in files:
        basename = os.path.basename(pathname)
        extension = get_file_extension(pathname)
        if extension in extensions or basename in extensions:
            outfiles.append(pathname)
    return outfiles


def main(args):
    """Identifies and verifies files that should have the desired boilerplate.

    Retrieves the lists of files to be validated and tests each one in turn.
    If all files contain correct boilerplate, this function terminates
    normally. Otherwise it prints the name of each non-conforming file and
    exists with a non-zero status code.
    """
    refs = get_refs(args)
    filenames = get_files(refs.keys(), args)
    nonconforming_files = []
    for filename in filenames:
        if not has_valid_header(filename, refs):
            nonconforming_files.append(filename)
    if nonconforming_files:
        print('%d files have incorrect boilerplate headers:' % len(
            nonconforming_files))
        for filename in sorted(nonconforming_files):
            print(os.path.relpath(filename, args.rootdir))
        sys.exit(1)


if __name__ == "__main__":
    ARGS = get_args()
    main(ARGS)
