// Fills benchoparity -ref with bencho.dev captures. Run with ego-browser, prefixing the
// settings it reads, e.g.
//   (echo 'globalThis.OUT="/abs/ref"; globalThis.SPACE=6; globalThis.START=0; globalThis.END=184;'; cat fetch-refs.mjs) | ego-browser nodejs
// OUT/_slugs.json lists the slugs in bencho.dev order (bencho.Catalog order).
const START = Number(globalThis.START ?? 0), END = Number(globalThis.END ?? 999);
const OUT = globalThis.OUT;
const fs = await import("node:fs/promises");
const task = await taskSpace(globalThis.SPACE ?? "bencho reference capture");
const page = task.page("p1");
const slugs = JSON.parse(await fs.readFile(`${OUT}/_slugs.json`, "utf8"));
const done = [];
for (let i = START; i < Math.min(END, slugs.length); i++) {
  const slug = slugs[i];
  const tag = String(i + 1).padStart(3, "0") + "-" + slug;
  try {
    await page.goto(`https://bencho.dev/blocks/${slug}`, { timeout: 30000 });
    await page.waitForSelector(".dtl-stage", { timeout: 15000 });
    await page.waitForTimeout(1800);
    const meta = await page.evaluate(() => {
      const st = document.querySelector(".dtl-stage").getBoundingClientRect();
      const side = document.querySelector(".dtl-side");
      const info = document.querySelector("[class*=dtl-info]");
      return { stage: { x: st.x, y: st.y, w: st.width, h: st.height }, side: side ? side.innerText : "", info: info ? info.innerText : "" };
    });
    const clip = { x: meta.stage.x, y: meta.stage.y, width: meta.stage.w, height: meta.stage.h };
    const cx = clip.x + clip.width / 2, cy = clip.y + clip.height / 2;
    await page.mouse.move(clip.x + 5, clip.y + clip.height - 5);
    await page.waitForTimeout(300);
    await page.screenshot({ path: `${OUT}/${tag}-a-rest.png`, clip });
    await page.mouse.move(cx, cy, { label: "hover block center" });
    await page.waitForTimeout(900);
    await page.screenshot({ path: `${OUT}/${tag}-b-hover.png`, clip });
    await page.mouse.down();
    await page.waitForTimeout(700);
    await page.screenshot({ path: `${OUT}/${tag}-c-press.png`, clip });
    await page.mouse.up();
    await page.waitForTimeout(900);
    await page.screenshot({ path: `${OUT}/${tag}-d-after.png`, clip });
    await fs.writeFile(`${OUT}/${tag}.json`, JSON.stringify({ slug, ...meta }, null, 1));
    done.push(tag);
  } catch (e) {
    console.log("FAIL", tag, e.message.slice(0, 200));
  }
}
console.log("done", done.length, "of", Math.min(END, slugs.length) - START, "total slugs", slugs.length);
