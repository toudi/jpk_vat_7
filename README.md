jpk vat z deklaracją
====================
Program do konwersji pliku / plików csv na plik JPK oraz wysyłania go
na serwer ministerstwa.

Napisałem go, ponieważ ministerstwo finansów udostępniło nową specyfikację JPK, ale nie udostępniło wejściowych plików CSV które można wrzucić do ich konwertera. Ponadto, oficjalny program ministerstwa nie wspiera nowej wersji
schemy

Ważna uwaga - dla zachowania kompatybilności z plikiem z ministerstwa używam średników
a nie przecinków.

Druga ważna uwaga - wypełniamy tylko te kolumny których potrzebujemy, dzięki temu
nie ma potrzeby wrzucania dziesiątek pustych pól - większość pól w JPK i tak jest opcjonalna.


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

```
KodFormularza;KodFormularza.kodSystemowy
ABC;DEF
```

co spowoduje wygenerowanie struktury w postaci:

```
<tns:KodFormularza kodSystemowy="DEF">ABC</tns:KodFormularza>
```

Pojedynczy plik
---------------

W wersji z pojedyńczym plikiem wszystki wiersze parsowane są z pojedyńczego wejściowego
pliku CSV. Przykład (celowo daję tylko kilka kolumn):

```
KodFormularza;KodFormularza.kodSystemowy;Pouczenia
JPK;JPK_V7M (1);
;;1
```

Katalog
-------

Przy wywołaniu programu z nazwą katalogu program będzie szukał następujących plików do
parsowania:

naglowek.csv (informacje nagłówkowe)
deklaracja.csv (informacje formularza VAT-7)
podmiot.csv
sprzedaz.csv (wiersze sprzedaży oraz wiersz kontrolny sprzedaży)
zakup.csv (wiersze zakupu oraz wiersz kontrolny zakupu)

Uwaga: kolumny KodFormularzaDekl oraz WariantFormularzaDekl nadpisywane są w sekcji nagłówka a nie w sekcji deklaracja.

Opis sekcji:
------------
Nagłówek - mapowany na gałąź tns:Naglowek
Kolumna startowa: KodFormularza

Podmiot - mapowany na gałąź tns:Podmiot > tns:OsobaFizyczna lub tns:Podmiot > tns:OsobaNiefizyczna
Kolumna startowa: typPodmiotu. Jeśli kolumna ma wartość "F" zostanie wygenerowana struktura tns:OsobaFizyczna; Jeśli kolumna ma wartość "NF" zostanie wygenerowana
struktura tns:OsobaNiefizyczna

Deklaracja - mapowany na gałąź tns:Deklaracja > tns:PozycjeSzczegolowe
Kolumna startowa: Pouczenia

Wiersze Sprzedaży - mapowane na gałęzie tns:SprzedazWiersz
Kolumna startowa: LpSprzedazy

Wiersz kontrolny sprzedaży - mapowany na gałąź tns:SprzedazCtrl
Kolumna startowa: LiczbaWierszySprzedazy

Wiersze kupna - mapowane na gałęzie tns:ZakupWiersz
Kolumna startowa: LpZakupu

Wiersz kontrolny zakupu - mapowany na gałąź tns:ZakupCtrl
Kolumna startowa: LiczbaWierszyZakupow

kompilacja:
-----------

go build