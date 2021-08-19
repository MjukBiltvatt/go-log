# go-log
Detta paket är ett standardiserat sätt att föra loggar på när det gäller program som är skrivna i Go. Det är detta paket som i framtiden kommer att skriva loggar på ett sätt som tillåter statistik att räknas ut i ett separat program.

## Användning
För att man ska kunna börja logga i sitt program behöver man först skapa en logg. Detta görs med hjälp av en utav två funktioner i detta paket. Dessa funktioner är beskrivna nedan tillsammans med vad det är som skiljer dem åt. 

### För text-loggar
För loggar som enkelt kan läsas av en människa bör funktionen `NewTextLogger` användas. Den kräver två strängar som parametrar varav den första representerar den katalog som loggen bör skapas i och den andra representerar loggens filnamn. Skulle katalogen för loggen inte existera när `NewTextLogger` körs så kommer den att skapas på disken innan loggfilen. Skulle något gå fel när katalogen ska skapas så returneras ett error. Ett exempel för användning ges nedan.
```go
main, err := log.NewTextLogger("my_log_dir", "main.log")
if err != nil {
    handleErr(err)
}
```

### För JSON-loggar
De loggar som bör skickas ut till externa system bör skrivas av i en logg av typen JSON. Att skapa en sådan logg är lika enkelt som för ovanstående, den enda skillnaden är att funktionen `NewJSONLogger` åkallas istället för `NewTextLogger`. Två parametrar krävs, katalognamnet och filnamnet för loggen. Om katalogen inte existerar när funktionen blir åkallad så kommer den att skapas. Skulle detta resultera i ett error så kommer detta att returneras till klienten. Ett exempel på hur man skapar en JSON logg presenteras nedan. 
```go
main, err := log.NewJSONLogger("my_log_dir", "main.log")
if err != nil {
    handleErr(err)
}
```

### Loggar
Det finns tre nivåer av loggar tillgängligt med detta paket, dessa är beskrivna nedan tillsammans med hur man skapar en logg av den typen.

#### Info 
En info-logg skapas genom att kalla på `Info`-funktionen som är bunden till `Log` gränssnittet. Funktionen tar en sträng som paramater tillsammans med en godtycklig mängd fält som skapas med funktionen `Field`. Ett exempel för att skapa en info logg är givet nedan.
```go
main, err := log.NewTextLogger("my_log_dir", "main.log")
if err != nil {
    handleErr(err)
}
main.Info("this is an info log")
main.Info("this info log has fields", log.Field("rating", 5))
```

#### Error
En error-logg skapas genom att kalla på `Log`-gränsnittets `Error` metod. Funktionen tar samma parametrar som `Info` och de behandlas på samma sätt. Den enda skillnaden är alltså loggens grad som istället blir `ERROR`. Utöver vanliga fält kan funktionen `ErrField` användas för att skriva ut det error som uppstod på ett standardiserat sätt. Exempel på detta är givet nedan.
```go
main, err := log.NewTextLogger("my_log_dir", "main.log")
if err != nil {
    handleErr(err)
}
main.Error("this is an error")
main.Error("this error log has fields", log.Field("rating", 5))
if err := somethingRisky(); err != nil {
    main.Error("could not perform risky operation", log.ErrField(err))
    handleErr(err)
}
```
Det är även möjligt att hämta mängden errors som har loggats för en given logg med metoden `Errors`. Denna metod returnerar ett heltal som representerar mängden errors som har loggats vid den tidpunkt då metoden åkallas. Ett enkelt exempel för dess användning är givet nedan.
```go
main, err := log.NewTextLogger("my_log_dir", "main.log")
if err != nil {
    handleErr(err)
}
main.Error("this is an error")
main.Error("this error log has fields", log.Field("rating", 5))
main.Info("several errors has occured", log.Field("amount", main.Errors()))
```

#### Warning
Loggning av varningar sker precis som ovantående loggtyper, den enda skillnaden är att metoden `Warn` används istället för `Info` eller `Error`. Likt `Error` så kan även antalet varningsloggar som har skapats hämtas med hjälp av metoden `Warnings`. Exempel för att skapa varningsloggar och hämta dess mängd är givet nedan.
```go
main, err := log.NewTextLogger("my_log_dir", "main.log")
if err != nil {
    handleErr(err)
}
main.Warn("this is a warning")
main.Warn("this warning log has fields", log.Field("rating", 5))
main.Info("several warnings has occured ", log.Field("amount", main.Warnings()))
```

### Ihopkopplade loggar
Det är möjligt att skapa flera loggar och koppla dessa till varandra i en parent-child relation. Detta betyder att alla loggar som skapas av föräldren kommer även att loggas av alla dess barn. Om en logg skapas specifikt på ett barn så kommer förälderloggen dock inte att föra samma logg. För att koppla ihop två loggar kan metoden `Attach` användas. Ett exempel för detta är givet nedan.
```go
main, err := log.NewTextLogger("my_log_dir", "main.log")
if err != nil {
    handleErr(err)
}
detailed, err := log.NewTextLogger("my_log_dir", "detailed.log")
if err != nil {
    handleErr(err)
}

main.Attach(detailed)
main.Info("this log will be in both main.log and detailed.log")
detailed.Info("this log will only be written to detailed.log, not main.log")
```

### Viktigt
För att tvinga loggarna att skrivas till fil så bör `Flush` bli åkallat vid slutet av programmets körtid eller när loggen inte längre kommer att användas. När `Flush` körs på en förälderlogg så kallas även metoden på alla loggens barn. Ett exempel av använding är givet nedan.
```go
main, err := log.NewTextLogger("my_log_dir", "main.log")
if err != nil {
    handleErr(err)
}
defer main.Flush()
// ... continue program and log things
```
`Flush` kan returnera ett error, för att hantera detta med `defer` så kan följande användas.
```go
main, err := log.NewTextLogger("my_log_dir", "main.log")
if err != nil {
    handleErr(err)
}
defer func() {
    if err := main.Flush(); err != nil {
        handleErr(err)
    }
}()
```