package pmongo

//// MarshalJSONPB marshals ObjectId to JSONPB string
//func (o *ObjectId) MarshalJSONPB(m *jsonpb.Marshaler) ([]byte, error) {
//
//
//
//
//
//
//	s, err := m.MarshalToString(&wrappers.StringValue{Value: o.Value})
//	if err != nil {
//		return nil, err
//	}
//	return []byte(s), nil
//}
//
//// UnmarshalJSONPB unmarshal JSONPB string to ObjectId
//func (o *ObjectId) UnmarshalJSONPB(m *jsonpb.Unmarshaler, data []byte) error {
//	var id wrappers.StringValue
//	if err := m.Unmarshal(bytes.NewReader(data), &id); err != nil {
//		return err
//	}
//	o.Value = id.Value
//	return nil
//}
