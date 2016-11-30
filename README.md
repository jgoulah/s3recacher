# S3 ReCacher

S3 ReCacher is a program that allow AWS S3 Administrators to easily set (or reset) the Cache-Control HTTP header on all objects within a bucket.

## Installation

Download the code:

    go get github.com/johnvilsack/s3recacher

Then run s3recacher from the bin/ directory!

## Usage

S3 Recacher requires Go1.3 to be installed.
~~~    
Usage: s3recacher [options...] -AWS_ACCESS_KEY_ID="YOURKEYID" -AWS_SECRET_ACCESS_KEY="YOURSECRET" -bucket="YOURBUCKETID"

Options:
  -age  Set the max-age (in seconds). Default is 0.
  
  -maxobjs Stop after parsing this many S3 objects. Off by default.

  -startid Filename to begin with. Useful for restarting.
     S3 always returns requests in alphabetical order.
  
  -stopid Filename to end with. Useful for distributing load.
~~~

The process will take time. Due to the linear access required to manipulate S3 data, there doesn't appear to be an easy way to add concurrency to speed up the process.

## Example

    ./s3recacher --bucket=primarydotcom --prefix=images/home/update --age=315360000

## License

The MIT License (MIT)

Copyright (c) 2014 John Vilsack

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
