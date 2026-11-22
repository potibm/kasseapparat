import React from "react";
import { splitByMatches } from "../../utils/searchHighlight";

/** Shows a text with the parts matching the search in bold and underlined. */
const HighlightedText: React.FC<{ text: string; query: string }> = ({
  text,
  query,
}) => {
  // Each part starts at a different offset, so offset + text is a stable key.
  let offset = 0;

  return (
    <span>
      {splitByMatches(text, query).map((part) => {
        const key = `${offset}-${part.text}`;
        offset += part.text.length;

        return part.match ? (
          <span key={key} className="font-bold underline">
            {part.text}
          </span>
        ) : (
          <React.Fragment key={key}>{part.text}</React.Fragment>
        );
      })}
    </span>
  );
};

export default HighlightedText;
