# http://commoncrawl.org/the-data/get-started/
# https://gist.github.com/Smerity/56bc6f21a8adec920ebf

import gzip
import json
import requests
from io import BytesIO
import magic


mimeExtensions = {"text/html":".html", "text/plain":".txt"}

def getIndex(crawl, url):
    """
    Make request to commoncrawl index API to gather all offsets that contain pointed URL
    Arguments: 
        crawl: Crawl database which should be used, e.g 'CC-MAIN-2019-22';
        url: URL of site, offsets and other info of which should be returned.
    Returns a list of JSON objects with information about each file offset and other data.
    """
    resp = requests.get('http://index.commoncrawl.org/{0}-index?url={1}&output=json'.format(crawl, url))
    pages = [json.loads(x) for x in resp.content.strip().split('\n'.encode("utf-8"))]
    return pages

def saveContent(pages, saveTo):
    """
    Saves pages or text that were found in Common Crawl to choosen folder
        pages: info about found web pages from `getIndex function`
        saveTo: destination fodler, where save fetched web data
    """
    crawlStorage = 'https://commoncrawl.s3.amazonaws.com/'

    for i, page in enumerate(pages):
        offset, length = int(page['offset']), int(page['length'])
        offsetEnd = offset + length - 1  
        resp = requests.get(crawlStorage + page['filename'], headers={'Range': 'bytes={}-{}'.format(offset, offsetEnd)})
    
        rawData = BytesIO(resp.content)
        f = gzip.GzipFile(fileobj=rawData)

        data = f.read()
        warc, header, response = data.strip().split('\r\n\r\n'.encode("utf-8"), 2)
    
        mime = magic.from_buffer(response, mime=True)
        ext = mimeExtensions[mime]
        startURL = warc.find(b'WARC-Target-URI:') + 17
        endURL = warc.find(b'\r\nWARC-Payload-Digest')
        url = warc[startURL:endURL].decode("utf-8")
        urlClear = "".join(["%"+str(ord(c)) if c in ["/", "\\", ":", "?"] else c for c in url])
    
        with open("{0}/{1}{2}".format(saveTo, urlClear, ext), "wb+") as f:
            f.write(response)
        print("Processing [{0}]: {1}/{2}".format(url, i+1, len(pages)))


if __name__ == "__main__":
    pages = getIndex(crawl="CC-MAIN-2019-22", url="example.com")
    print(pages)
    #saveContent(pages = pages, saveTo = "data")
