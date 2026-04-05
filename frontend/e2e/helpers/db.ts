import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export async function resetDatabase() {
  const databaseDir = path.resolve(__dirname, "../../../backend/data/");
  const templateDbPath = path.resolve(databaseDir, "e2e-clean.db");
  const activeTestDbPath = path.resolve(databaseDir, "e2e-work.db");

  try {
    await new Promise((resolve) => setTimeout(resolve, 50));

    fs.copyFileSync(templateDbPath, activeTestDbPath);
    // eslint-disable-next-line no-console
    console.log("🔄 SQLite database reset successful.");
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error("❌ Error while resetting database:", error);
    throw error;
  }
}
