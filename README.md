jpk_gen
=======
Program do konwersji pliku / plików csv na plik JPK.

Napisałem go, ponieważ ministerstwo finansów udostępniło nową specyfikację JPK, ale nie udostępniło wejściowych plików CSV które można wrzucić do ich konwertera.

Instrukcja obsługi
------------------
Program może pracować w dwóch trybach: parsowania pojedyńczego pliku lub katalogu z 
plikami. Dodatkowo, załączam przykładowe pliki wejściowe

Sekcje
------
Program zakłada, że plik (pliki) wejściowe CSV składają się z sekcji. Każda sekcja 
rozpoczyna się od kolumny a następnie kolejne kolumny parsowane są do sekcji. Niektóre 
pozycje w pliku JPK wymagają atrybutów. Wówczas nagłówek musi być w postaci

kolumna.atrybut

na przykład:

KodFormularza;KodFormularza.kodSystemowy
ABC;DEF

co spowoduje wygenerowanie struktury w postaci:

<tns:KodFormularza kodSystemowy="DEF">ABC</tns:KodFormularza>

Kolejność kolumn ma znaczenie w tym sensie, że kiedy program znajdzie kolejną sekcję
to kolumny dopisywane są do tej sekcji. puste  kolumny są ignorowane, co oznacza, że
technicznie wykonalne jest stworzenie pliku będącego kompatybilnie wstecznym z tym od
ministerstwa, np:

KodFormularza; .... (kolumny); Rok
JPK;;;; ( ;;;; ); 2019
;;;;;1;0;0;;;;;

itd.

Jak jednak widać, takie podejście jest kiepsko skalowalne

Pojedynczy plik
~~~~~~~~~~~~~~~
Zachowałem ten tryb dla kompatybilności z plikiem dostarczonym przez ministerstwo,
ale moim skromnym zdaniem używanie tej wersji mija się z celem. aby plik był kompatybilny
wstecz, jakiekolwiek nowe kolumny muszą być dodawane na końcu co oznacza, że kolejne 
wiersze muszą mieć gigantyczne ilości średników aby być sparsowane. Dodatkowo, nie
wszystkie kolumny są obowiązkowe więc tym bardziej mamy dużą ilość średników

Katalog
~~~~~~~

Przy wywołaniu programu z nazwą katalogu program będzie szukał następujących plików do
parsowania:

naglowek.csv (informacje nagłówkowe)
deklaracja.csv (informacje formularza VAT-7)
podmiot.csv
sprzedaz.csv (wiersze sprzedaży oraz wiersz kontrolny sprzedaży)
zakup.csv (wiersze zakupu oraz wiersz kontrolny zakupu)

Uwaga: kolumny KodFormularzaDekl oraz WariantFormularzaDekl nadpisywane są w sekcji nagłówka a nie w sekcji deklaracja.

Opis sekcji:
~~~~~~~~~~~~
Nagłówek
Kolumna startowa: KodFormularza

Deklaracja
Kolumna startowa: Pouczenia

Wiersze Sprzedaży
Kolumna startowa: LpSprzedazy

Wiersz kontrolny sprzedaży
Kolumna startowa: LiczbaWierszySprzedazy

Wiersze kupna
Kolumna startowa: LpZakupu

Wiersz kontrolny zakupu
Kolumna startowa: LiczbaWierszyZakupow

kompilacja:
-----------

go build