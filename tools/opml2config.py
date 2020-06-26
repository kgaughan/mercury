#!/usr/bin/env python3

import sys
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
    escapes = str.maketrans(ESCAPES)

    first = True
    for path in sys.argv[1:]:
        doc = minidom.parse(path)
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
