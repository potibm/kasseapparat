import React from "react";
import { Avatar } from "flowbite-react";

export const GuestlistAvatar: React.FC<{ name: string }> = ({ name }) => {
  const getInitials = (name: string) => {
    // Remove all non-alphabetical characters except spaces
    const cleanedName = name.replaceAll(/[^a-zA-Z\s]/g, "").trim();

    // Split the cleaned name into words
    const words = cleanedName.split(" ").filter((word) => word.length > 0);

    // If there's only one word, take the first letter twice
    if (words.length === 1) {
      return words[0][0].toUpperCase();
    }

    // For multiple words, take the first letter of the first and last word
    const firstInitial = words[0][0].toUpperCase();
    const lastInitial = words[words.length - 1][0].toUpperCase();

    return firstInitial + lastInitial;
  };

  const initials = getInitials(name);

  return <Avatar placeholderInitials={initials} size="md" rounded />;
};

export default GuestlistAvatar;
