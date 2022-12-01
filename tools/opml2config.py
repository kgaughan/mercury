#!/usr/bin/env python3

from os import path
import sys
from typing import Iterator, Tuple
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


def extract_feeds(doc: minidom.Document) -> Iterator[Tuple[str, str]]:
    for node in doc.getElementsByTagName("outline"):
        if node.getAttribute("type") == "rss":
            if feed := node.getAttribute("xmlUrl"):
                yield (feed, node.getAttribute("text"))


def main():
    if len(sys.argv) == 2 and sys.argv[1] == "-h":
        print("usage:", path.basename(sys.argv[0]), "<opml>*", file=sys.stderr)
        sys.exit(0)

    escapes = str.maketrans(ESCAPES)
    for filename in sys.argv[1:]:
        if not path.exists(filename):
            sys.exit(f"error: no such file: {filename}")
        for feed, name in extract_feeds(minidom.parse(filename)):
            print("[[feed]]")
            print(f'name = "{name.translate(escapes)}"')
            print(f'feed = "{feed.translate(escapes)}"')
            print()


if __name__ == "__main__":
    main()
