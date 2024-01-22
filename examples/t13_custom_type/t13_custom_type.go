package main

// my custom type
type Vector4D struct {
   w float64
   x float64
   y float64
   z float64
};

// add this just in case, if it is necessary to register it with
// Groot2 publisher.
// You will need to add `RegisterJsonDefinition<Vector4D>(ToJson);` in you main
func  ToJson( dest map[string]interface{},  pose Vector4D){
  dest["w"] = pose.w;
  dest["x"] = pose.x;
  dest["y"] = pose.y;
  dest["z"] = pose.z;
}

namespace BT {

template <> inline
    Vector4D convertFromString(StringView key)
{
  const auto parts = BT::splitString(key, ',');
  if (parts.size() != 4)
  {
    throw BT::RuntimeError("invalid input)");
  }

  Vector4D output;
  output.w = convertFromString<double>(parts[0]);
  output.x     = convertFromString<double>(parts[1]);
  output.y     = convertFromString<double>(parts[2]);
  output.z = convertFromString<double>(parts[3]);
  return output;
}

}
