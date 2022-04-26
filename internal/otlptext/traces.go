package otlptext

import (
	"go.opentelemetry.io/collector/model/pdata"
)

// NewTextTracesMarshaler returns a serializer.TracesMarshaler to encode to OTLP text bytes.
func NewTextTracesMarshaler() pdata.TracesMarshaler {
	return textTracesMarshaler{}
}

type textTracesMarshaler struct{}

// MarshalTraces pdata.Traces to OTLP text.
func (textTracesMarshaler) MarshalTraces(td pdata.Traces) ([]byte, error) {
	buf := dataBuffer{}
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		buf.logEntry("ResourceSpans #%d", i)
		rs := rss.At(i)
		buf.logAttributeMap("Resource labels", rs.Resource().Attributes())
		ilss := rs.InstrumentationLibrarySpans()
		for j := 0; j < ilss.Len(); j++ {
			buf.logEntry("InstrumentationLibrarySpans #%d", j)
			ils := ilss.At(j)
			buf.logInstrumentationLibrary(ils.InstrumentationLibrary())

			spans := ils.Spans()
			for k := 0; k < spans.Len(); k++ {
				buf.logEntry("Span #%d", k)
				span := spans.At(k)
				buf.logAttr("Trace ID", span.TraceID().HexString())
				buf.logAttr("Parent ID", span.ParentSpanID().HexString())
				buf.logAttr("ID", span.SpanID().HexString())
				buf.logAttr("Name", span.Name())
				buf.logAttr("Kind", span.Kind().String())
				buf.logAttr("Start time", span.StartTimestamp().String())
				buf.logAttr("End time", span.EndTimestamp().String())

				buf.logAttr("Status code", span.Status().Code().String())
				buf.logAttr("Status message", span.Status().Message())

				buf.logAttributeMap("Attributes", span.Attributes())
				buf.logEvents("Events", span.Events())
				buf.logLinks("Links", span.Links())
			}
		}
	}

	return buf.buf.Bytes(), nil
}
