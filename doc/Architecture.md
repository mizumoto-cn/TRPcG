# How TRPG is built

## Header

So now lets talk about the `Header`, which will be defined in [/header/](/header/).

As for header definition, we separate request and respond header into two files : [req](header/req_header.go) and [res](header/res_header.go).

To reuse created Request/Response Header objects, TRPcG implements buffer pools. When a header finishes its job, TRPcG will reset its status using a ResetHeader() method, then they'll be thrown back into the pool again.

## IO

The IO operating methods are built in [/codec/io.go](/codec/io.go) .

-> write to IO stream: `sendFrame()`

Firstly put len(data) into stream, if it isn't 0, write []byte into it.