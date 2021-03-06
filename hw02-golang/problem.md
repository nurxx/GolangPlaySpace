Имате следния интерфейс:
```
type Task interface {
    Execute(int) (int, error)
}
```
(Трябва да го включите в предаденото от вас домашно.)

Имплементирайте следните функции (конструктори), които да приемат и връщат задачи от горния тип:

1. Функция `Pipeline(tasks ...Task) Task` със следните свойства:
- Приема произволен брой Task обекти и връща написан от вас тип, който също имплементира Task интерфейса.
- При извикването на метода `Execute()` на върнатия от `Pipeline()` обект, методите `Execute()` на всички задачи от подадените tasks трябва да се изпълнят последователно.
- За аргумент на първата задача от tasks се използва аргумента на Execute(), а за аргумент на всяка следваща се използва резултата от предишната.
- Като краен резултат на `Execute()` метода на pipeline обекта се връща резултата от `Execute()` на последната задача от tasks.
- Ако няма подадени задачи на `Pipeline()`, `Execute()` трябва да върне грешка. Ако някоя от задачите върне грешка, целият `Execute()` на pipeline-а приключва (не се изпълняват повече задачи) и връща грешка.
Ето как изглежда това в код. Нека си направим следния прост тип adder, който не хапе, а събира цели числа до `127` и имплементира интерфейса `Task`:
  ```
  type adder struct {
      augend int
  }

  func (a adder) Execute(addend int) (int, error) {
      result := a.augend + addend
      if result > 127 {
          return 0, fmt.Errorf("Result %d exceeds the adder threshold", a)
      }
      return result, nil
  }
  ```
Ето и как очакваме да се държи върнатия от `Pipeline()` тип:

```
if res, err := Pipeline(adder{50}, adder{60}).Execute(10); err != nil {
    fmt.Printf("The pipeline returned an error\n")
} else {
    fmt.Printf("The pipeline returned %d\n", res)
}
```
Това би трябвало да изведе "The pipeline returned 120". Но ако имахме `Pipeline(adder{20}, adder{10}, adder{-50}).Execute(100)`, би трябвало да получим на екрана `"The pipeline returned an error".`

2. Функция `Fastest(tasks ...Task) Task` със следните свойства:
- Отново приема произволен брой Task обекти и връща написан от вас тип, който също имплементира Task интерфейса.
- При извикването на метода `Execute()` на върнатия от `Fastest()` обект, методите `Execute()` на всички задачи от подадените tasks трябва да се изпълнят конкурентно и да се върне резултата (или грешката) на тази задача, която завърши първа.
- Като аргумент на всички задачи от tasks се подава едно и също число - аргументът, с който е извикан `Execute()` на върнатия от `Fastest()` обект.
- Ако няма подадени задачи на `Fastest()`, `Execute()` трябва да върне грешка.
- Постарайте се да не оставяте "висящи" горутини, ще смъкваме точки.
Ето прост пример, преизползвайки adder от горния пример:
  ```type lazyAdder struct {
      adder
      delay time.Duration
  }

  func (la lazyAdder) Execute(addend int) (int, error) {
      time.Sleep(la.delay * time.Millisecond)
      return la.adder.Execute(addend)
  }
  ```
би трябвало да получим `42` от следния код:
```
f := Fastest(
    lazyAdder{adder{20}, 500},
    lazyAdder{adder{50}, 300},
    adder{41},
)
f.Execute(1)
```
3. Функция `Timed(task Task, timeout time.Duration) Task` със следните свойства:
- Приема една задача от тип `Task` и `timeout` време и връща написан от вас тип, който също имплементира `Task` интерфейса.
- При извикването на метода `Execute()` на връщания обект, изпълнява `task.Execute()` със същата стойност и връща получения резултат или грешка ако задачата приключи в зададеното от timeout време. - Ако не успее да приключи за това време, връща грешка.
- Постарайте се да не оставяте "висящи" горутини, ще смъкваме точки.
Ето пример, преизползвайки `lazyAdder`:
  _, e1 := Timed(lazyAdder{adder{20}, 50}, 2*time.Millisecond).Execute(2)
  r2, e2 := Timed(lazyAdder{adder{20}, 50}, 300*time.Millisecond).Execute(2)
Очакваме първия ред да върне грешка (т.е. `e1 != nil`), a вторият ред да е ок и `r2` да съдържа резултата `22`.

4. Функция `ConcurrentMapReduce(reduce func(results []int) int, tasks ...Task) Task` със следните свойства:
- Приема `reduce` функция и произволен брой `Task` обекти и връща написан от вас тип, който също имплементира `Task` интерфейса.
- При извикването на метода `Execute()` на върнатия от `ConcurrentMapReduce()` обект, методите `Execute()` на всички задачи от подадените tasks трябва да се изпълнят конкурентнo.
- Като аргумент на всички задачи се подава аргументът, с който е извикан `Execute()`.
- Ако няма подадени задачи на `ConcurrentMapReduce()`, `Execute()` трябва да върне грешка. 
- Ако някоя от функциите fail-не, `Execute()` трябва веднага да върне грешка. Ако всички задачи приключат успешно, трябва да се извика reduce с техните резултати (в произволен ред) и резултатът от reduce да бъде върнат като резултат на функцията.
- Постарайте се да не оставяте "висящи" горутини, ще смъкваме точки.
Следният код:
```
  reduce := func(results []int) int {
      smallest := 128
      for _, v := range results {
          if v < smallest {
              smallest = v
          }
      }
      return smallest
  }

  mr := ConcurrentMapReduce(reduce, adder{30}, adder{50}, adder{20})
      if res, err := mr.Execute(5); err != nil {
      fmt.Printf("We got an error!\n")
  } else {
      fmt.Printf("The ConcurrentMapReduce returned %d\n", res)
  }
```
би трябвало да изведе `The ConcurrentMapReduce returned 25`.

5. Функция `GreatestSearcher(errorLimit int, tasks <-chan Task) Task` със следните свойства:
- Приема максимален допустим брой на грешките `errorLimit` и небуфериран канал за четене tasks, по който асинхронно могат да ѝ се подават задачи за изпълнение. Отново връща написан от вас тип, който също имплементира `Task` интерфейса.
- При извикването на `Execute()` от върнатия `Task` трябва всичките задачи от канала tasks да започнат greedily да се изпълняват конкурентно. Искаме да няма блокиране, щом ние пуснем задача по този канал, вашия тип трябва да я прочете от канала и да извика нейния `Execute()` метод
- `Execute()` метода на задачата трябва да приключи след като ние затворим tasks канала и всички вече подадени задачи от него са приключили.
- Като резултат `Execute()` метода на вашия тип, след приключването на всики задачи, трябва да върне най-голямото число, което някоя задача е върнала. Но ако повече от errorLimit задачи са върнали грешка или по tasks не бъдат подадени никакви задачи, `Execute()` трябва да върне грешка.
Пример:
  ```
  tasks := make(chan Task)
  gs := GreatestSearcher(2, tasks) // Приемаме 2 грешки

  go func() {
      tasks <- adder{4}
      tasks <- lazyAdder{adder{22}, 20}
      tasks <- adder{125} // Това е първата "допустима" грешка (защото 125+10 > 127)
      time.Sleep(50 * time.Millisecond)
      tasks <- adder{32} // Това би трябвало да "спечели"

      // Това би трябвало да timeout-не и да е втората "допустима" грешка
      tasks <- Timed(lazyAdder{adder{100}, 2000}, 20*time.Millisecond)

      // Ако разкоментираме това, gs.Execute() трябва да върне грешка
      // tasks <- adder{127} // трета (и недопустима) грешка
      close(tasks)
  }()
  result, err := gs.Execute(10)
  ```
Очакваме да получим 42 като result. Но ако разкоментираме реда с `tasks <- adder{127}`, тогава gs.`Execute()` трябва да върне грешка (т.е. `err != nil`).