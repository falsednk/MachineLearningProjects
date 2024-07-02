package preprocessing

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

// Загрузить датафрейм
func GetDataframe(path string) dataframe.DataFrame {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(file)
	return df
}

// Посчитать пустые строки
func CountIsNan(s series.Series) int {
	count_isnan := 0
	for i := 0; i < s.Len(); i++ {
		if s.Elem(i).IsNA() {
			count_isnan++
		}
	}
	return count_isnan
}

// Получить уникальные значения
func GetUniqueValues(s series.Series) []any {
	uniqueMap := make(map[any]bool)
	var uniqueValues []any
	switch s.Type() {
	case "string":
		for i := 0; i < s.Len(); i++ {
			elem := s.Elem(i).String()
			if _, ok := uniqueMap[elem]; !ok {
				uniqueMap[elem] = true
				if !s.Elem(i).IsNA() {
					uniqueValues = append(uniqueValues, elem)
				}
			}
		}
	case "float":
		for i := 0; i < s.Len(); i++ {
			elem := s.Elem(i).Float()
			if _, ok := uniqueMap[elem]; !ok {
				uniqueMap[elem] = true
				if !s.Elem(i).IsNA() {
					uniqueValues = append(uniqueValues, elem)
				}
			}
		}
	case "int":
		for i := 0; i < s.Len(); i++ {
			elem, _ := s.Elem(i).Int()
			if _, ok := uniqueMap[elem]; !ok {
				uniqueMap[elem] = true
				if !s.Elem(i).IsNA() {
					uniqueValues = append(uniqueValues, elem)
				}
			}
		}
	}
	return uniqueValues
}

// Выгрузить общую информацию о датафрейме
func GetInfoDf(df dataframe.DataFrame) {
	fmt.Printf("Size dataframe: rows %d, columns %d\n", df.Nrow(), df.Ncol())
	fmt.Println("_ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ \n")
	for i, name := range df.Names() {
		fmt.Println("Column: ", name, "  |Type: ", df.Types()[i], "|IsNanSum: ", CountIsNan(df.Col(name)), "|Count unique value: ", len(GetUniqueValues(df.Col(name))))
	}
}

// Получить часть даты
func GetPartDate(s series.Series, part string) series.Series {
	dateLayout := "2006-01-02"
	datePart := make([]int, s.Len())
	for i := 0; i < s.Len(); i++ {
		date, err := time.Parse(dateLayout, s.Elem(i).String())
		year, week := date.ISOWeek()
		if err != nil {
			datePart[i] = 0
		}
		if part == "Year" {
			datePart[i] = year
		} else if part == "Month" {
			datePart[i] = int(date.Month())
		} else if part == "Day" {
			datePart[i] = int(date.Day())
		} else {
			datePart[i] = week
		}
	}
	return series.New(datePart, series.Int, part)
}

// Добавить расширенные данные по дате
func DateApply(df dataframe.DataFrame) dataframe.DataFrame {
	df = df.Mutate(GetPartDate(df.Col("Date"), "Year"))
	df = df.Mutate(GetPartDate(df.Col("Date"), "Month"))
	df = df.Mutate(GetPartDate(df.Col("Date"), "Week"))
	df = df.Mutate(GetPartDate(df.Col("Date"), "Day"))
	return df
}

// Получить индексы пропусков по условию, заданному другой колонкой
func GetIndices(s, condition series.Series, value any) []int {
	var indices []int
	for i := 0; i < s.Len(); i++ {
		var elem any
		switch value.(type) {
		case string:
			val := condition.Elem(i).String()
			elem = val
		case float64:
			val := condition.Elem(i).Float()
			elem = val
		case int:
			val, _ := condition.Elem(i).Int()
			elem = val
		case bool:
			val, _ := condition.Elem(i).Bool()
			elem = val
		}
		if (s.Elem(i).IsNA()) && (elem == value) {
			indices = append(indices, i)
		}
	}
	return indices
}

// Заменить пропуски. Только для float64
func FillNa(s series.Series, value float64, column string) series.Series {
	newSeries := make([]float64, s.Len())
	for i := 0; i < s.Len(); i++ {
		if s.Elem(i).IsNA() {
			newSeries[i] = value
		} else {
			newSeries[i] = s.Elem(i).Float()
		}
	}
	return series.New(newSeries, series.Float, column)
}

// Реализация OneHotEncoding
func OneHotEncoding(df dataframe.DataFrame, s series.Series) dataframe.DataFrame {

	uniqueValues := GetUniqueValues(s)
	for part, value := range uniqueValues {
		var newCol []int
		var newName string = value.(string)
		for i := 0; i < s.Len(); i++ {
			elem := s.Elem(i).String()
			if elem == value {
				newCol = append(newCol, 1)
			} else {
				newCol = append(newCol, 0)
			}
		}
		if part != 0 {
			df = df.Mutate(series.New(newCol, series.Int, newName))
		}
	}
	return df
}

// Получить колонки по типу
func GetNameColumnsByType(df dataframe.DataFrame, typeColumn string) []string {

	var nameColumns []string
	for i, name := range df.Names() {
		if string(df.Types()[i]) == typeColumn {
			nameColumns = append(nameColumns, name)
		}
	}
	return nameColumns
}
