package plugin

import (
	"bytes"
	"github.com/M09ic/go-ntlmssp"
	"github.com/chainreactors/gogo/v2/pkg"
	"github.com/chainreactors/utils/iutils"
)

var prelogin = []byte{
	0x12, 0x01, 0x00, 0x58, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x06, 0x01, 0x00, 0x25,
	0x00, 0x01, 0x02, 0x00, 0x26, 0x00, 0x01, 0x03, 0x00, 0x27, 0x00, 0x04, 0x04, 0x00, 0x2b, 0x00,
	0x01, 0x05, 0x00, 0x2c, 0x00, 0x24, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

var sspiMessage = []byte{
	0x10, 0x01, 0x01, 0xb3, 0x00, 0x00, 0x01, 0x00, 0xab, 0x01, 0x00, 0x00, 0x04, 0x00, 0x00, 0x74,
	0x40, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x2a, 0x2a, 0x2a, 0x2a, 0x00, 0x00, 0x00, 0x00,
	0xe0, 0x83, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5e, 0x00, 0x09, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x70, 0x00, 0x21, 0x00, 0xb2, 0x00, 0x0e, 0x00,
	0xce, 0x00, 0x04, 0x00, 0xd2, 0x00, 0x21, 0x00, 0x14, 0x01, 0x00, 0x00, 0x14, 0x01, 0x07, 0x00,
	0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x22, 0x01, 0x7e, 0x00, 0xa0, 0x01, 0x00, 0x00, 0xa0, 0x01,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x41, 0x00, 0x4e, 0x00, 0x4f, 0x00, 0x4e, 0x00, 0x59, 0x00,
	0x4d, 0x00, 0x4f, 0x00, 0x55, 0x00, 0x53, 0x00, 0x43, 0x00, 0x6f, 0x00, 0x72, 0x00, 0x65, 0x00,
	0x20, 0x00, 0x2e, 0x00, 0x4e, 0x00, 0x65, 0x00, 0x74, 0x00, 0x20, 0x00, 0x53, 0x00, 0x71, 0x00,
	0x6c, 0x00, 0x43, 0x00, 0x6c, 0x00, 0x69, 0x00, 0x65, 0x00, 0x6e, 0x00, 0x74, 0x00, 0x20, 0x00,
	0x44, 0x00, 0x61, 0x00, 0x74, 0x00, 0x61, 0x00, 0x20, 0x00, 0x50, 0x00, 0x72, 0x00, 0x6f, 0x00,
	0x76, 0x00, 0x69, 0x00, 0x64, 0x00, 0x65, 0x00, 0x72, 0x00, 0x31, 0x00, 0x30, 0x00, 0x2e, 0x00,
	0x32, 0x00, 0x30, 0x00, 0x30, 0x00, 0x2e, 0x00, 0x32, 0x00, 0x31, 0x00, 0x35, 0x00, 0x2e, 0x00,
	0x31, 0x00, 0x30, 0x00, 0x38, 0x00, 0xa0, 0x01, 0x00, 0x00, 0x43, 0x00, 0x6f, 0x00, 0x72, 0x00,
	0x65, 0x00, 0x20, 0x00, 0x2e, 0x00, 0x4e, 0x00, 0x65, 0x00, 0x74, 0x00, 0x20, 0x00, 0x53, 0x00,
	0x71, 0x00, 0x6c, 0x00, 0x43, 0x00, 0x6c, 0x00, 0x69, 0x00, 0x65, 0x00, 0x6e, 0x00, 0x74, 0x00,
	0x20, 0x00, 0x44, 0x00, 0x61, 0x00, 0x74, 0x00, 0x61, 0x00, 0x20, 0x00, 0x50, 0x00, 0x72, 0x00,
	0x6f, 0x00, 0x76, 0x00, 0x69, 0x00, 0x64, 0x00, 0x65, 0x00, 0x72, 0x00, 0x54, 0x00, 0x64, 0x00,
	0x73, 0x00, 0x54, 0x00, 0x65, 0x00, 0x73, 0x00, 0x74, 0x00, 0x60, 0x7c, 0x06, 0x06, 0x2b, 0x06,
	0x01, 0x05, 0x05, 0x02, 0xa0, 0x72, 0x30, 0x70, 0xa0, 0x30, 0x30, 0x2e, 0x06, 0x0a, 0x2b, 0x06,
	0x01, 0x04, 0x01, 0x82, 0x37, 0x02, 0x02, 0x0a, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x82, 0xf7, 0x12,
	0x01, 0x02, 0x02, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x12, 0x01, 0x02, 0x02, 0x06, 0x0a,
	0x2b, 0x06, 0x01, 0x04, 0x01, 0x82, 0x37, 0x02, 0x02, 0x1e, 0xa2, 0x3c, 0x04, 0x3a, 0x4e, 0x54,
	0x4c, 0x4d, 0x53, 0x53, 0x50, 0x00, 0x01, 0x00, 0x00, 0x00, 0xb7, 0xb2, 0x08, 0xe2, 0x09, 0x00,
	0x09, 0x00, 0x31, 0x00, 0x00, 0x00, 0x09, 0x00, 0x09, 0x00, 0x28, 0x00, 0x00, 0x00, 0x0a, 0x00,
	0x61, 0x4a, 0x00, 0x00, 0x00, 0x0f, 0x41, 0x4e, 0x4f, 0x4e, 0x59, 0x4d, 0x4f, 0x55, 0x53, 0x57,
	0x4f, 0x52, 0x4b, 0x47, 0x52, 0x4f, 0x55, 0x50, 0x01, 0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00,
	0x00, 0x00, 0xff,
}

func mssqlScan(result *pkg.Result) {
	result.Port = "1433"
	target := result.GetTarget()
	conn, err := pkg.NewSocket("tcp", target, RunOpt.Delay)
	if err != nil {
		return
	}
	_, err = conn.Request(prelogin, 4096)

	if err != nil {
		return
	}

	r2, err := conn.Request(sspiMessage, 4096)
	off_ntlm := bytes.Index(r2, []byte("NTLMSSP"))
	if off_ntlm <= 0 {
		return
	}
	data := r2[off_ntlm:]
	defer conn.Close()
	result.Open = true
	result.Status = "mssql"
	result.Protocol = "mssql"
	result.AddNTLMInfo(iutils.ToStringMap(ntlmssp.NTLMInfo(data)), "mssql")
}
