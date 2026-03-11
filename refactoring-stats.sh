#!/bin/bash

# Ordner festlegen (Standard: src)
TARGET_DIR=${1:-frontend/src}

echo "--------------------------------------------------"
echo "📊 Migrations-Status für: $TARGET_DIR"
echo "--------------------------------------------------"

# Hilfsfunktion zum Zählen
count_stats() {
    local ext_pattern=$1
    local files=$(find "$TARGET_DIR" -type f -name "$ext_pattern" | wc -l)
    local loc=$(find "$TARGET_DIR" -type f -name "$ext_pattern" -exec cat {} + | wc -l)
    echo "$files $loc"
}

# Daten sammeln
read js_files js_loc <<< $(count_stats "*.js*")
read ts_files ts_loc <<< $(count_stats "*.ts*")

# Berechnungen
total_files=$((js_files + ts_files))
total_loc=$((js_loc + ts_loc))

# Prozentwerte (mit awk für Fließkommazahlen)
if [ "$total_files" -gt 0 ]; then
    pct_files=$(awk "BEGIN {printf \"%.2f\", ($ts_files/$total_files)*100}")
    pct_loc=$(awk "BEGIN {printf \"%.2f\", ($ts_loc/$total_loc)*100}")
else
    pct_files=0
    pct_loc=0
fi

# Ausgabe
printf "%-15s | %-10s | %-10s\n" "Typ" "Dateien" "LOC"
echo "--------------------------------------------------"
printf "%-15s | %-10d | %-10d\n" "JS/JSX (Alt)" "$js_files" "$js_loc"
printf "%-15s | %-10d | %-10d\n" "TS/TSX (Neu)" "$ts_files" "$ts_loc"
echo "--------------------------------------------------"
printf "%-15s | %-10d | %-10d\n" "Gesamt" "$total_files" "$total_loc"
echo ""
echo "🚀 Umstellungsgrad (Dateien): $pct_files %"
echo "📈 Umstellungsgrad (LOC):     $pct_loc %"
echo "--------------------------------------------------"