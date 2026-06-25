// Package scannercore is the JNI bridge surface for the Scan-Zero Go core.
//
// Serialized Primitive Protocol:
//   - Inputs:    a plain string, or a JSON-encoded array of absolute paths.
//   - Responses: a single primitive (string), or a thrown error.
//     (An exported func whose last return is error becomes a Java method
//     that throws.)
//   - All file IO is confined to the caller-supplied cache directory.
package scannercore

func Ping(input string) string {

	return "pong: " + input
}

func unpackPaths(input string) ([]string, error) {
	//string input goes to json.Unmarshal
	//json.Unmarshal goes into a []string, and handle error cases.
	var a []string
	var err error
	return a, err
}
