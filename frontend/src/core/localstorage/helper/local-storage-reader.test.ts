import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { 
  getCurrentReaderId, 
  setCurrentReaderId, 
  clearCurrentReaderId, 
  COOKIE_NAME 
} from './local-storage-reader';
import { createLogger } from '@core/logger/logger';

vi.mock('@core/logger/logger', () => {
  const mockLogger = {
    error: vi.fn(),
  };
  return {
    createLogger: vi.fn(() => mockLogger),
  };
});

const VALID_READER_ID = 'rdr_ABCDEFGHIJKLMNOPQRSTUVWXYZ'; 
const VALID_READER_ID_2 = 'rdr_12345678901234567890123456';
const INVALID_READER_ID = 'rdr_invalid_lowercase';

describe('Reader ID Management', () => {
  const mockLogger = createLogger('Config');

  beforeEach(() => {
    document.cookie = `${COOKIE_NAME}=; Max-Age=0; path=/`;
    vi.clearAllMocks();
  });

  afterEach(() => {
     document.cookie = `${COOKIE_NAME}=; Max-Age=0; path=/`;
  });

  describe('getCurrentReaderId', () => {
    it('should be undefined when no cookie exists', () => {
      const result = getCurrentReaderId();
      expect(result).toBeUndefined();
    });

    it('should return the Reader ID when a valid cookie exists', () => {
      document.cookie = `${COOKIE_NAME}=${VALID_READER_ID}`;
      const result = getCurrentReaderId();
      expect(result).toBe(VALID_READER_ID);
    });

    it('should also parse URL-encoded cookies correctly', () => {
      document.cookie = `${COOKIE_NAME}=${encodeURIComponent(VALID_READER_ID_2)}`;
      const result = getCurrentReaderId();
      expect(result).toBe(VALID_READER_ID_2);
    });

    it('should be undefined when the cookie value is invalid (Zod fails)', () => {
      document.cookie = `${COOKIE_NAME}=${INVALID_READER_ID}`;
      const result = getCurrentReaderId();
      expect(result).toBeUndefined();
    });

    it('should be undefined when the cookie has invalid URI encoding', () => {
      // a '%' without following characters throws an error in decodeURIComponent
      document.cookie = `${COOKIE_NAME}=rdr_DEFECTIVE%`; 
      const result = getCurrentReaderId();
      expect(result).toBeUndefined();
    });
  });

  describe('setCurrentReaderId', () => {
    it('should set the cookie when the Reader ID is valid', () => {
      setCurrentReaderId(VALID_READER_ID);
      
      expect(document.cookie).toContain(`${COOKIE_NAME}=${VALID_READER_ID}`);
      expect(mockLogger.error).not.toHaveBeenCalled();
    });

    it('should not set the cookie and log an error when the Reader ID is invalid', () => {
      setCurrentReaderId(INVALID_READER_ID);
      
      expect(document.cookie).not.toContain(COOKIE_NAME);
      
      expect(mockLogger.error).toHaveBeenCalledTimes(1);
      expect(mockLogger.error).toHaveBeenCalledWith(
        "Invalid Reader ID provided",
        expect.any(Function) 
      );
    });
  });

  describe('clearCurrentReaderId', () => {
    it('should clear the cookie (set Max-Age=0)', () => {
      document.cookie = `${COOKIE_NAME}=${VALID_READER_ID}`;
      expect(document.cookie).toContain(COOKIE_NAME);

      clearCurrentReaderId();
      
      expect(document.cookie).not.toContain(COOKIE_NAME);
    });
  });
});