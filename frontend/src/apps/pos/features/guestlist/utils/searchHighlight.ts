/** Letters that do not decompose into a base letter plus a diacritic. */
const SPECIAL_FOLDS: Record<string, string> = {
  ß: "ss",
  æ: "ae",
  œ: "oe",
  ø: "o",
  ł: "l",
};

/**
 * Lower-cases one character and strips its diacritics, like the backend search does
 * (ADR / PR #988): "Ü" → "u", "ß" → "ss".
 */
export const foldChar = (char: string): string => {
  const lower = char.toLowerCase();

  return (
    SPECIAL_FOLDS[lower] ??
    lower.normalize("NFD").replaceAll(/\p{Diacritic}/gu, "")
  );
};

export const foldText = (text: string): string =>
  Array.from(text).map(foldChar).join("");

export interface TextPart {
  text: string;
  match: boolean;
}

/**
 * Splits a text into the parts that match a search and the parts in between, ignoring
 * case and diacritics, so the search "MUL" highlights "Mül" in "Müller".
 */
export const splitByMatches = (text: string, query: string): TextPart[] => {
  const foldedQuery = foldText(query.trim());
  if (!text || foldedQuery === "") return text ? [{ text, match: false }] : [];

  // Fold character by character and remember where each folded character came
  // from, because folding can change the length ("ß" → "ss").
  const chars = Array.from(text);
  let folded = "";
  const origin: number[] = [];
  chars.forEach((char, index) => {
    const foldedChar = foldChar(char);
    folded += foldedChar;
    for (let i = 0; i < foldedChar.length; i++) origin.push(index);
  });

  const matched = new Array<boolean>(chars.length).fill(false);
  let from = folded.indexOf(foldedQuery);
  while (from !== -1) {
    for (let i = from; i < from + foldedQuery.length; i++) {
      matched[origin[i]] = true;
    }
    from = folded.indexOf(foldedQuery, from + foldedQuery.length);
  }

  const parts: TextPart[] = [];
  chars.forEach((char, index) => {
    const last = parts.at(-1);
    if (last && last.match === matched[index]) {
      last.text += char;
    } else {
      parts.push({ text: char, match: matched[index] });
    }
  });

  return parts;
};
