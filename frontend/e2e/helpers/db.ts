import fs from "node:fs";
import path from "node:path";

export function resetDatabase() {
  const databaseDir = path.resolve(process.cwd(), "../backend/data/");
  const templateDbPath = path.resolve(databaseDir, "e2e-clean.db");
  const activeTestDbPath = path.resolve(databaseDir, "e2e-work.db");

  try {
    fs.copyFileSync(templateDbPath, activeTestDbPath);
    // eslint-disable-next-line no-console
    console.log("🔄 SQLite database reset successful.");
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error("❌ Error while resetting database:", error);
    throw error;
  }
}
