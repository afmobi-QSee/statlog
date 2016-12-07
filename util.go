package statlog

import "bytes"

func (this *StatusUdp)BufferStringJoin(args ...string) string{
	var buffer bytes.Buffer
	buffer.WriteString(this.StatEnv)

	for _, args := range args {
		buffer.WriteString(args)
	}
	return buffer.String()
}
