package zgrab2

var readBufSize = 1024

// StreamReadEnvelope - Reads remote scan requests from a stream.
// func StreamReadEnvelope(connection io.ReadWriteCloser) (req *Request, err error) {
//
//         // Read the first four bytes to determine data length
//         dataLengthBuf := make([]byte, 4) // Size of uint32
//         _, err = connection.Read(dataLengthBuf)
//         if err != nil {
//                 return nil, fmt.Errorf("(read msg-length): %v", err)
//         }
//         dataLength := int(binary.LittleEndian.Uint32(dataLengthBuf))
//
//         // Read the length of the data, keep in mind each call to .Read() may not
//         // fill the entire buffer length that we specify, so instead we use two buffers
//         // readBuf is the result of each .Read() operation, which is then concatinated
//         // onto dataBuf which contains all of data read so far and we keep calling
//         // .Read() until the running total is equal to the length of the message that
//         // we're expecting or we get an error.
//         // readBuf := make([]byte, readBufSize)
//         dataBuf := make([]byte, 0)
//         totalRead := 0
//         for {
//                 // Compute the precise length of the temporary buffer
//                 var readBuf []byte
//                 if dataLength-len(dataBuf) > readBufSize {
//                         readBuf = make([]byte, readBufSize)
//                 } else {
//                         readBuf = make([]byte, (dataLength - len(dataBuf)))
//                 }
//
//                 n, err := connection.Read(readBuf)
//                 dataBuf = append(dataBuf, readBuf[:n]...)
//                 totalRead += n
//                 if totalRead == dataLength {
//                         break
//                 }
//                 if err != nil {
//                         break
//                 }
//         }
//         if err != nil {
//                 return req, err
//         }
//
//         // Unmarshal the protobuf envelope
//         req = &RemoteRequest{}
//         err = json.Unmarshal(dataBuf, req)
//         if err != nil {
//                 return nil, fmt.Errorf("unmarshaling request error: %v", err)
//         }
//         return req, nil
// }
