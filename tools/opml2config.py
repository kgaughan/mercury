#!/usr/bin/env python3

import sys
from os import path
from xml.dom import minidom

ESCAPES = {
    '"': r"\"",
    "\\": r"\\",
    "\b": r"\b",
    "\f": r"\f",
    "\n": r"\n",
    "\r": r"\r",
    "\t": r"\t",
}


def main():
    if len(sys.argv) == 2 and sys.argv[1] == "-h":
        print("usage:", path.basename(sys.argv[0]), "<opml>", file=sys.stderr)
        sys.exit(0)

    escapes = str.maketrans(ESCAPES)

    first = True
    for filename in sys.argv[1:]:
        if not path.exists(filename):
            print("error: no such file:", filename, file=sys.stderr)
            sys.exit(1)
        doc = minidom.parse(filename)
        for node in doc.getElementsByTagName("outline"):
            node_type = node.getAttribute("type")
            if node_type != "rss":
                continue
            if first:
                first = False
            else:
                print()
            name = node.getAttribute("text").translate(escapes)
            feed = node.getAttribute("xmlUrl").translate(escapes)
            if feed:
                print("[[feed]]")
                print(f'name = "{name}"')
                print(f'feed = "{feed}"')


if __name__ == "__main__":
    main()
