// Fills benchoparity's -ref directory from bencho.dev: <slug>.jpg, the 2x still
// bencho.dev shows for each block, and specs.json, each block's display size,
// whether its preview is live, its description, props and controls. Run with
// ego-browser, prefixing the directory and the slugs (bencho.Catalog order):
//
//	(echo 'globalThis.OUT="/abs/ref"; globalThis.SLUGS=["thinking-orb","dice"];'; cat fetch-refs.mjs) | ego-browser nodejs
const fs = await import("node:fs/promises");
const { OUT, SLUGS } = globalThis;
await fs.mkdir(OUT, { recursive: true });
const task = await taskSpace("bencho fetch-refs");
const page = task.page("p1");
const specs = {};
for (const slug of SLUGS) {
  try {
    await page.goto(`https://bencho.dev/blocks/${slug}`, { timeout: 30000 });
    await page.waitForSelector("[class*=dtl-info]", { timeout: 15000 });
    specs[slug] = await page.evaluate(() => {
      const sc = document.querySelector(".dtl-stage .dtl-scale")?.getBoundingClientRect();
      return {
        live: !document.querySelector(".dtl-stage .pro-still"),
        size: sc ? [Math.round(sc.width), Math.round(sc.height)] : null,
        side: document.querySelector(".dtl-side")?.innerText ?? "",
        info: document.querySelector("[class*=dtl-info]")?.innerText ?? "",
      };
    });
    const res = await fetch(`https://bencho.dev/stills/${slug}.jpg`);
    if (res.ok) await fs.writeFile(`${OUT}/${slug}.jpg`, Buffer.from(await res.arrayBuffer()));
  } catch (e) {
    console.log("FAIL", slug, e.message.slice(0, 160));
  }
}
await fs.writeFile(`${OUT}/specs.json`, JSON.stringify(specs, null, 1));
console.log(`fetch-refs: ${Object.keys(specs).length} of ${SLUGS.length} blocks in ${OUT}`);
await task.finish({ keep: [] });
