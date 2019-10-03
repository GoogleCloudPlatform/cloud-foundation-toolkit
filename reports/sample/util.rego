package validator.gcp.lib

# has_field returns whether an object has a field
has_field(object, field) {
        object[field]
}

# False is a tricky special case, as false responses would create an undefined document unless
# they are explicitly tested for
has_field(object, field) {
        object[field] == false
}

has_field(object, field) = false {
        not object[field]
        not object[field] == false
}

# get_default returns the value of an object's field or the provided default value.
# It avoids creating an undefined state when trying to access an object attribute that does
# not exist
get_default(object, field, _default) = output {
        has_field(object, field)
        output = object[field]
}

get_default(object, field, _default) = output {
        has_field(object, field) == false
        output = _default
}

bool_to_str(bool_value) = "true" {
    bool_value
}

bool_to_str(bool_value) = "false" {
    not bool_value
}

is_null_str(value) = "null" {
    count({value} & {{},[]}) == 1
}

is_null_str(value) = "defined" {
    count({value} & {{},[]}) == 0
}