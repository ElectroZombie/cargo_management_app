package database

import (
    "database/sql"
    "errors"
    "fmt"
    "net/mail"
    "regexp"
    "strings"
    "time"
)

var (
    ErrNotFound = errors.New("record not found")
    ErrValidation = errors.New("validation failed")
    ErrConflict = errors.New("record conflicts with existing data")
)

type ValidationError struct { Field string; Message string }
func (e ValidationError) Error() string { return fmt.Sprintf("%s: %s", e.Field, e.Message) }
func invalid(field, message string) error { return fmt.Errorf("%w: %s", ErrValidation, ValidationError{field, message}) }
func required(value, field string) error { if strings.TrimSpace(value)=="" { return invalid(field,"is required") }; return nil }
func maxLen(value, field string, max int) error { if len([]rune(strings.TrimSpace(value)))>max{return invalid(field,fmt.Sprintf("must be at most %d characters",max))};return nil }
func positive(value float64, field string) error { if value<0{return invalid(field,"must be zero or greater")};return nil }
func positiveID(value int, field string) error { if value<=0{return invalid(field,"must be greater than zero")};return nil }
func email(value, field string) error { if strings.TrimSpace(value)=="" { return nil }; if _,err:=mail.ParseAddress(value);err!=nil{return invalid(field,"must be a valid email address")};return nil }
func date(value, field string, requiredDate bool) error { value=strings.TrimSpace(value);if value==""&&!requiredDate{return nil};if value==""{return invalid(field,"is required")};if _,err:=time.Parse("2006-01-02",value);err!=nil{return invalid(field,"must use YYYY-MM-DD format")};return nil }
func oneOf(value, field string, allowed ...string) error { for _,v:=range allowed{if value==v{return nil}};return invalid(field,"has an unsupported value") }
func validateDriver(v DriverDTO) error { if e:=required(v.Name,"name");e!=nil{return e};if e:=required(v.LicenseNumber,"license_number");e!=nil{return e};if e:=maxLen(v.Name,"name",200);e!=nil{return e};if e:=maxLen(v.LicenseNumber,"license_number",100);e!=nil{return e};return email(v.Email,"email") }
func validateLoader(v LoaderDTO) error { if e:=required(v.Name,"name");e!=nil{return e};if e:=maxLen(v.Name,"name",200);e!=nil{return e};return email(v.Email,"email") }
func validateVehicle(v VehicleDTO) error { for _,x:=range[][2]string{{v.VecID,"vec_id",""},{v.LicenseID,"license_id",""},{v.VIN,"vin",""},{v.LicenseExpiration,"license_expiration",""}}{if e:=required(x[0],x[1]);e!=nil{return e}};if e:=date(v.LicenseExpiration,"license_expiration",true);e!=nil{return e};if e:=positiveID(v.DriverID,"driver_id");e!=nil{return e};if v.TrailerNumber<0{return invalid("trailer_number","must be zero or greater")};return nil}
func validateWell(v WellDTO) error { if e:=required(v.Name,"name");e!=nil{return e};return required(v.Client,"client") }
func validateLoad(v LoadDTO) error { for _,x:=range[][2]string{{v.Date,"date",""},{v.LoadNumber,"load_number",""},{v.TicketNumber,"ticket_number",""},{v.VecID,"vec_id",""}}{if e:=required(x[0],x[1]);e!=nil{return e}};if e:=date(v.Date,"date",true);e!=nil{return e};for _,x:=range[][2]float64{{v.Miles,"miles",""},{v.NetWeight,"net_weight",""},{v.Tons,"tons",""},{v.TonRate,"ton_rate",""},{v.TotalValue,"total_value",""}}{if e:=positive(x[0],x[1]);e!=nil{return e}};for _,x:=range[][2]int{{v.LoaderID,"loader_id",""},{v.WellID,"well_id",""},{v.DriverID,"driver_id",""}}{if e:=positiveID(x[0],x[1]);e!=nil{return e}};if v.TrailerNumber<0{return invalid("trailer_number","must be zero or greater")};return nil}

func affected(result sql.Result, entity string) error { n,err:=result.RowsAffected();if err!=nil{return fmt.Errorf("%s update: %w",entity,err)};if n==0{return fmt.Errorf("%w: %s",ErrNotFound,entity)};return nil }
func databaseError(operation string, err error) error { if err==nil{return nil};message:=strings.ToLower(err.Error());if strings.Contains(message,"unique constraint")||strings.Contains(message,"constraint failed")&&strings.Contains(message,"unique"){return fmt.Errorf("%w during %s: %v",ErrConflict,operation,err)};if strings.Contains(message,"foreign key constraint"){return fmt.Errorf("%w during %s: referenced record does not exist",ErrValidation,operation)};return fmt.Errorf("%s: %w",operation,err) }
func scanDriver(row interface{Scan(...any)error})(*Driver,error){var v Driver;err:=row.Scan(&v.ID,&v.Name,&v.LicenseNumber,&v.Phone,&v.Email,&v.Address,&v.CreatedAt,&v.UpdatedAt);return &v,err}
func scanLoader(row interface{Scan(...any)error})(*Loader,error){var v Loader;err:=row.Scan(&v.ID,&v.Name,&v.Phone,&v.Email,&v.Address,&v.Company,&v.CreatedAt,&v.UpdatedAt);return &v,err}
func scanVehicle(row interface{Scan(...any)error})(*Vehicle,error){var v Vehicle;err:=row.Scan(&v.ID,&v.VecID,&v.LicenseID,&v.TrailerNumber,&v.VIN,&v.LicenseExpiration,&v.DriverID,&v.OverweightPermitID,&v.CreatedAt,&v.UpdatedAt);return &v,err}
func scanWell(row interface{Scan(...any)error})(*Well,error){var v Well;err:=row.Scan(&v.ID,&v.Name,&v.Client,&v.Location,&v.CreatedAt,&v.UpdatedAt);return &v,err}
func scanLoad(row interface{Scan(...any)error})(*Load,error){var v Load;err:=row.Scan(&v.ID,&v.Date,&v.LoadNumber,&v.LoaderID,&v.WellID,&v.TicketNumber,&v.Miles,&v.DriverID,&v.VecID,&v.TrailerNumber,&v.NetWeight,&v.Tons,&v.TonRate,&v.TotalValue,&v.Status,&v.TicketID,&v.CreatedAt,&v.UpdatedAt);return &v,err}

var identifierPattern=regexp.MustCompile(`^[A-Za-z0-9._/-]+$`)
