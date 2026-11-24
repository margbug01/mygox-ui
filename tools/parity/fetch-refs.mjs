// Screenshots every particle card on coss.com/ui/particles into OUT, one
// <id>.png each, for tools/parity. Run with: ego-browser nodejs < fetch-refs.mjs
// after writing REGISTRY from coss's apps/ui/registry/registry-particles.ts:
//   node --experimental-strip-types -e 'import("./registry-particles.ts").then(m => require("fs").writeFileSync("particles.json", JSON.stringify(m.particles)))'
// Particles already in OUT are skipped, so a rerun fills only what failed.
const REGISTRY = "/absolute/path/particles.json";
const OUT = "/absolute/path/coss-ref";

const fs = await import("node:fs");
const up = JSON.parse(fs.readFileSync(REGISTRY, "utf8"));
fs.mkdirSync(OUT, { recursive: true });
const done = new Set(fs.readdirSync(OUT).map((f) => f.replace(".png", "")));
const task = await taskSpace("coss particle references");
const page = task.page("p1");
const cards = "div.relative.flex.min-w-0";
const missed = [];
for (const cat of [...new Set(up.flatMap((p) => p.categories))]) {
  const want = up.filter((p) => p.categories.includes(cat) && !done.has(p.name));
  if (!want.length) continue;
  await page.goto("https://coss.com/ui/particles?tags=" + encodeURIComponent(cat), { timeout: 60000 });
  await page.waitForTimeout(2500);
  for (let k = 0; k < 12; k++) {
    await page.evaluate(() => window.scrollBy(0, 1200));
    await page.waitForTimeout(300);
  }
  // A card's description is the last .truncate span: its footer's.
  const descs = await page.evaluate((sel) => [...document.querySelectorAll(sel)].map((c) => [...c.querySelectorAll("span.truncate")].pop()?.textContent || ""), cards);
  const inCat = up.filter((p) => p.categories.includes(cat));
  for (const p of want) {
    const k = inCat.filter((q) => q.description === p.description).indexOf(p);
    const idx = descs.map((d, i) => (d === p.description ? i : -1)).filter((i) => i >= 0)[k];
    if (idx === undefined) continue;
    await page.evaluate(({ sel, i }) => document.querySelectorAll(sel)[i].scrollIntoView({ block: "center" }), { sel: cards, i: idx });
    await page.waitForTimeout(250);
    const clip = await page.evaluate(({ sel, i }) => {
      const b = document.querySelectorAll(sel)[i].getBoundingClientRect();
      return { x: b.x + scrollX, y: b.y + scrollY, width: b.width, height: b.height };
    }, { sel: cards, i: idx });
    await page.screenshot({ path: OUT + "/" + p.name + ".png", clip });
    done.add(p.name);
  }
}
for (const p of up) if (!done.has(p.name)) missed.push(p.name);
await task.finish({ keep: [] });
console.log({ have: done.size, missed });
