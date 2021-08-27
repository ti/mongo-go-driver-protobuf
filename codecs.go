package codecs

import (
	"errors"
	"reflect"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"google.golang.org/protobuf/types/known/timestamppb"
	
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ti/mongo-go-driver-protobuf/pmongo"
)

var (
	// Protobuf wrappers types
	boolValueType   = reflect.TypeOf(wrapperspb.BoolValue{})
	bytesValueType  = reflect.TypeOf(wrapperspb.BytesValue{})
	doubleValueType = reflect.TypeOf(wrapperspb.DoubleValue{})
	floatValueType  = reflect.TypeOf(wrapperspb.FloatValue{})
	int32ValueType  = reflect.TypeOf(wrapperspb.Int32Value{})
	int64ValueType  = reflect.TypeOf(wrapperspb.Int64Value{})
	stringValueType = reflect.TypeOf(wrapperspb.StringValue{})
	uint32ValueType = reflect.TypeOf(wrapperspb.UInt32Value{})
	uint64ValueType = reflect.TypeOf(wrapperspb.UInt64Value{})

	// Protobuf Timestamp type
	timestampType = reflect.TypeOf(timestamppb.Timestamp{})

	// Time type
	timeType = reflect.TypeOf(time.Time{})

	// ObjectId type
	objectIDType          = reflect.TypeOf(pmongo.ObjectId{})
	objectIDPrimitiveType = reflect.TypeOf(primitive.ObjectID{})

	// Codecs
	wrapperValueCodecRef = &wrapperValueCodec{}
	timestampCodecRef    = &timestampCodec{}
	objectIDCodecRef     = &objectIDCodec{}
)

// wrapperValueCodec is codec for Protobuf type wrappers
type wrapperValueCodec struct {
}

// EncodeValue encodes Protobuf type wrapper value to BSON value
func (e *wrapperValueCodec) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	val = val.FieldByName("Value")
	enc, err := ectx.LookupEncoder(val.Type())
	if err != nil {
		return err
	}
	return enc.EncodeValue(ectx, vw, val)
}

// DecodeValue decodes BSON value to Protobuf type wrapper value
func (e *wrapperValueCodec) DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	val = val.FieldByName("Value")
	enc, err := ectx.LookupDecoder(val.Type())
	if err != nil {
		return err
	}
	return enc.DecodeValue(ectx, vr, val)
}

// timestampCodec is codec for Protobuf Timestamp
type timestampCodec struct {
}

// EncodeValue encodes Protobuf Timestamp value to BSON value
func (e *timestampCodec) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	v, ok := val.Interface().(timestamppb.Timestamp)
	if !ok {
		return errors.New("value is not timestamppb.Timestamp")
	}
	t := v.AsTime()
	enc, err := ectx.LookupEncoder(timeType)
	if err != nil {
		return err
	}
	return enc.EncodeValue(ectx, vw, reflect.ValueOf(t.In(time.UTC)))
}

// DecodeValue decodes BSON value to Timestamp value
func (e *timestampCodec) DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	enc, err := ectx.LookupDecoder(timeType)
	if err != nil {
		return err
	}
	var t time.Time
	if err = enc.DecodeValue(ectx, vr, reflect.ValueOf(&t).Elem()); err != nil {
		return err
	}
	ts := timestamppb.New(t.In(time.UTC))
	val.Set(reflect.ValueOf(*ts))
	return nil
}

// objectIDCodec is codec for Protobuf ObjectId
type objectIDCodec struct {
}

// EncodeValue encodes Protobuf ObjectId value to BSON value
func (e *objectIDCodec) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	v, _ := val.Interface().(pmongo.ObjectId)
	// Create primitive.ObjectId from string
	id, err := primitive.ObjectIDFromHex(v.Value)
	if err != nil {
		return err
	}
	enc, err := ectx.LookupEncoder(objectIDPrimitiveType)
	if err != nil {
		return err
	}
	return enc.EncodeValue(ectx, vw, reflect.ValueOf(id))
}

// DecodeValue decodes BSON value to ObjectId value
func (e *objectIDCodec) DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	enc, err := ectx.LookupDecoder(objectIDPrimitiveType)
	if err != nil {
		return err
	}
	var id primitive.ObjectID
	if err = enc.DecodeValue(ectx, vr, reflect.ValueOf(&id).Elem()); err != nil {
		return err
	}
	oid := *pmongo.NewObjectId(id)
	if err != nil {
		return err
	}
	val.Set(reflect.ValueOf(oid))
	return nil
}

// Register registers Google protocol buffers types codecs
func Register(rb *bsoncodec.RegistryBuilder) *bsoncodec.RegistryBuilder {
	return rb.RegisterCodec(boolValueType, wrapperValueCodecRef).
		RegisterCodec(bytesValueType, wrapperValueCodecRef).
		RegisterCodec(doubleValueType, wrapperValueCodecRef).
		RegisterCodec(floatValueType, wrapperValueCodecRef).
		RegisterCodec(int32ValueType, wrapperValueCodecRef).
		RegisterCodec(int64ValueType, wrapperValueCodecRef).
		RegisterCodec(stringValueType, wrapperValueCodecRef).
		RegisterCodec(uint32ValueType, wrapperValueCodecRef).
		RegisterCodec(uint64ValueType, wrapperValueCodecRef).
		RegisterCodec(timestampType, timestampCodecRef).
		RegisterCodec(objectIDType, objectIDCodecRef)
}
