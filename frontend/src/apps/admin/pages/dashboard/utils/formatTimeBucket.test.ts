import { describe, it, expect } from "vitest";
import {
  parseTimeBucket,
  formatTimeBucketString,
  formatXAxisLabel,
  calculateGranularity,
  generateTimeBuckets,
  snapToGranularity,
} from "./formatTimeBucket";

describe("formatTimeBucket utilities", () => {
  describe("parseTimeBucket", () => {
    it("parses time bucket string to Date", () => {
      const date = parseTimeBucket("2024-01-15 14:30");
      expect(date.getUTCFullYear()).toBe(2024);
      expect(date.getUTCMonth()).toBe(0); // January is 0
      expect(date.getUTCDate()).toBe(15);
      expect(date.getUTCHours()).toBe(14);
      expect(date.getUTCMinutes()).toBe(30);
    });
  });

  describe("formatTimeBucketString", () => {
    it("formats date without minutes", () => {
      const date = new Date("2024-01-15T14:30:00Z");
      const result = formatTimeBucketString(date, false);
      expect(result).toBe("2024-01-15 14:00");
    });

    it("formats date with minutes", () => {
      const date = new Date("2024-01-15T14:30:00Z");
      const result = formatTimeBucketString(date, true);
      expect(result).toBe("2024-01-15 14:30");
    });
  });

  describe("formatXAxisLabel", () => {
    it("formats hourly label", () => {
      const result = formatXAxisLabel("2024-01-15 14:00", 60);
      expect(result).toMatch(/\d{2}-\d{2} \d{2}:00/);
    });

    it("formats sub-hourly label", () => {
      const result = formatXAxisLabel("2024-01-15 14:30", 15);
      expect(result).toMatch(/\d{2}-\d{2} \d{2}:30/);
    });

    it("returns empty string for empty input", () => {
      expect(formatXAxisLabel("", 60)).toBe("");
    });

    it("returns original value for invalid format", () => {
      expect(formatXAxisLabel("invalid", 60)).toBe("invalid");
    });
  });

  describe("calculateGranularity", () => {
    it("returns 15 minutes for span < 12 hours", () => {
      const min = new Date("2024-01-15T10:00:00Z");
      const max = new Date("2024-01-15T20:00:00Z");
      expect(calculateGranularity(min, max)).toBe(15);
    });

    it("returns 30 minutes for span 12-24 hours", () => {
      const min = new Date("2024-01-15T00:00:00Z");
      const max = new Date("2024-01-15T18:00:00Z");
      expect(calculateGranularity(min, max)).toBe(30);
    });

    it("returns 60 minutes for span > 24 hours", () => {
      const min = new Date("2024-01-15T00:00:00Z");
      const max = new Date("2024-01-17T00:00:00Z");
      expect(calculateGranularity(min, max)).toBe(60);
    });
  });

  describe("generateTimeBuckets", () => {
    it("generates hourly buckets", () => {
      const min = new Date("2024-01-15T10:00:00Z");
      const max = new Date("2024-01-15T12:00:00Z");
      const buckets = generateTimeBuckets(min, max, 60);
      expect(buckets).toHaveLength(3);
      expect(buckets[0]).toBe("2024-01-15 10:00");
      expect(buckets[1]).toBe("2024-01-15 11:00");
      expect(buckets[2]).toBe("2024-01-15 12:00");
    });

    it("generates 15-minute buckets", () => {
      const min = new Date("2024-01-15T10:00:00Z");
      const max = new Date("2024-01-15T11:00:00Z");
      const buckets = generateTimeBuckets(min, max, 15);
      expect(buckets).toHaveLength(5);
      expect(buckets[0]).toBe("2024-01-15 10:00");
      expect(buckets[1]).toBe("2024-01-15 10:15");
      expect(buckets[2]).toBe("2024-01-15 10:30");
      expect(buckets[3]).toBe("2024-01-15 10:45");
      expect(buckets[4]).toBe("2024-01-15 11:00");
    });
  });

  describe("snapToGranularity", () => {
    it("snaps to 15-minute granularity", () => {
      const date = new Date("2024-01-15T14:37:00Z");
      const snapped = snapToGranularity(date, 15);
      expect(snapped.getUTCMinutes()).toBe(30);
    });

    it("snaps to 30-minute granularity", () => {
      const date = new Date("2024-01-15T14:47:00Z");
      const snapped = snapToGranularity(date, 30);
      expect(snapped.getUTCMinutes()).toBe(30);
    });

    it("snaps to 60-minute granularity", () => {
      const date = new Date("2024-01-15T14:47:00Z");
      const snapped = snapToGranularity(date, 60);
      expect(snapped.getUTCMinutes()).toBe(0);
    });
  });
});
