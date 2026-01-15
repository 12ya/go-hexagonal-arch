package domain

// PortUpdate represents a partial update to a Port.
// Nil fields are not updated, non-nil fields are updated to the pointed value.
type PortUpdate struct {
	Name        *string
	Code        *string
	City        *string
	Country     *string
	Alias       *[]string
	Regions     *[]string
	Coordinates *[]float64
	Province    *string
	Timezone    *string
	Unlocs      *[]string
}
