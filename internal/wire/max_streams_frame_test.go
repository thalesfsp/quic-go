package wire

import (
	"bytes"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MAX_STREAMS frame", func() {
	Context("parsing", func() {
		It("accepts a frame for a bidirectional stream", func() {
			data := []byte{0x1c}
			data = append(data, encodeVarInt(0xdecaf)...)
			b := bytes.NewReader(data)
			f, err := parseMaxStreamsFrame(b, protocol.VersionWhatever)
			Expect(err).ToNot(HaveOccurred())
			Expect(f.Type).To(Equal(protocol.StreamTypeBidi))
			Expect(f.MaxStreams).To(BeEquivalentTo(0xdecaf))
			Expect(b.Len()).To(BeZero())
		})

		It("accepts a frame for a bidirectional stream", func() {
			data := []byte{0x1d}
			data = append(data, encodeVarInt(0xdecaf)...)
			b := bytes.NewReader(data)
			f, err := parseMaxStreamsFrame(b, protocol.VersionWhatever)
			Expect(err).ToNot(HaveOccurred())
			Expect(f.Type).To(Equal(protocol.StreamTypeUni))
			Expect(f.MaxStreams).To(BeEquivalentTo(0xdecaf))
			Expect(b.Len()).To(BeZero())
		})

		It("errors on EOFs", func() {
			data := []byte{0x1d}
			data = append(data, encodeVarInt(0xdeadbeefcafe13)...)
			_, err := parseMaxStreamsFrame(bytes.NewReader(data), protocol.VersionWhatever)
			Expect(err).NotTo(HaveOccurred())
			for i := range data {
				_, err := parseMaxStreamsFrame(bytes.NewReader(data[0:i]), protocol.VersionWhatever)
				Expect(err).To(HaveOccurred())
			}
		})
	})

	Context("writing", func() {
		It("for a bidirectional stream", func() {
			f := &MaxStreamsFrame{
				Type:       protocol.StreamTypeBidi,
				MaxStreams: 0xdeadbeef,
			}
			b := &bytes.Buffer{}
			Expect(f.Write(b, protocol.VersionWhatever)).To(Succeed())
			expected := []byte{0x1c}
			expected = append(expected, encodeVarInt(0xdeadbeef)...)
			Expect(b.Bytes()).To(Equal(expected))
		})

		It("for a unidirectional stream", func() {
			f := &MaxStreamsFrame{
				Type:       protocol.StreamTypeUni,
				MaxStreams: 0xdecafbad,
			}
			b := &bytes.Buffer{}
			Expect(f.Write(b, protocol.VersionWhatever)).To(Succeed())
			expected := []byte{0x1d}
			expected = append(expected, encodeVarInt(0xdecafbad)...)
			Expect(b.Bytes()).To(Equal(expected))
		})

		It("has the correct min length", func() {
			frame := MaxStreamsFrame{MaxStreams: 0x1337}
			Expect(frame.Length(protocol.VersionWhatever)).To(Equal(1 + utils.VarIntLen(0x1337)))
		})
	})
})