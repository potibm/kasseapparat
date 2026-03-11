#!/bin/bash

TARGET_DIR=${1:-frontend/src}

echo "-----------------------------------------------------------------------------------------------"
echo "🎯 Top-Kandidaten für die TS-Konvertierung"
echo "-----------------------------------------------------------------------------------------------"
printf "%-60s | %-8s | %-8s | %-10s\n" "File" "LOC" "Imports" "Score*"
echo "-----------------------------------------------------------------------------------------------"

# Alle .jsx und .js Dateien finden
find "$TARGET_DIR" -type f \( -name "*.jsx" -o -name "*.js" \) | while read -r file; do
    # Dateiname ohne Pfad für die Suche
    filename=$(basename "$file" | cut -d. -f1)
    
    # Wie oft wird diese Datei importiert? (Grob-Check via grep)
    import_count=$(grep -r "from .*$filename" "$TARGET_DIR" | grep -v "$file" | wc -l)
    
    # Zeilenanzahl
    loc=$(wc -l < "$file")
    
    # Score-Berechnung: (Imports * 10) + (LOC / 5)
    # Hohe Import-Zahl gewichtet stärker
    score=$(( (import_count * 10) + (loc / 5) ))
    
    printf "%-60s | %8d | %8d | %10d\n" "${file#$TARGET_DIR/}" "$loc" "$import_count" "$score"
done | sort -rn -k7 | head -n 15

echo "-----------------------------------------------------------------------------------------------"
echo "* Score = (Imports * 10) + (Zeilen / 5)"
echo "Tipp: Hoher Score = Hohe Priorität!"