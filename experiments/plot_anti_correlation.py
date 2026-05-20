"""Plot delta vs emb_d for the productionized recipe set across models.
Shows the universal anti-correlation finding."""
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
import numpy as np

# Per-model data: list of (recipe, delta, emb_d)
data = {
    "Sonnet 4.6": {
        "color": "#cc7722",
        "points": [
            ("R12-sigil", 0.327, 0.201),
            ("R12-anchor-duo-occult", 0.277, 0.191),
            ("R12-grimoire", 0.307, 0.238),
            ("R19-chord-anchor", 0.328, 0.251),
            ("R24-psalter-chord", 0.178, 0.303),
            ("psalter", 0.177, 0.327),
            ("R26-biblical-witness", 0.127, 0.332),
        ],
    },
    "Opus 4.7": {
        "color": "#5577dd",
        "points": [
            ("R12-sigil", 0.573, 0.244),
            ("R12-anchor-duo-occult", 0.596, 0.278),
            ("R12-grimoire", 0.411, 0.291),
            ("R19-chord-anchor", 0.516, 0.326),
            ("R24-psalter-chord", 0.388, 0.375),
            ("R26-biblical-witness", 0.503, 0.361),
        ],
    },
    "DeepSeek-pro v4": {
        "color": "#7733aa",
        "points": [
            ("R29-minimal", 0.272, 0.351),
            ("R12-sigil", 0.204, 0.398),
            ("envoy-extreme", 0.127, 0.367),
            ("R12-anchor-duo-occult", 0.022, 0.442),
            ("R12-grimoire", 0.018, 0.474),
            ("R19-chord-anchor", -0.085, 0.437),
            ("R24-psalter-chord", -0.104, 0.563),
            ("psalter", -0.118, 0.500),
            ("R26-biblical-witness", -0.158, 0.547),
        ],
    },
}


def pearson(xs, ys):
    n = len(xs)
    mx = sum(xs) / n
    my = sum(ys) / n
    cov = sum((x - mx) * (y - my) for x, y in zip(xs, ys)) / n
    vx = sum((x - mx) ** 2 for x in xs) / n
    vy = sum((y - my) ** 2 for y in ys) / n
    return cov / (vx**0.5 * vy**0.5)


fig, axes = plt.subplots(1, 3, figsize=(15, 5), sharey=False)

for ax, (model, info) in zip(axes, data.items()):
    pts = info["points"]
    embds = [p[2] for p in pts]
    deltas = [p[1] for p in pts]
    r = pearson(deltas, embds)
    ax.scatter(embds, deltas, c=info["color"], s=90, edgecolor="black", linewidth=0.5, zorder=3)
    for name, d, e in pts:
        ax.annotate(name, (e, d), xytext=(4, 4), textcoords="offset points",
                    fontsize=7, alpha=0.8)
    ax.axhline(0, color="gray", linewidth=0.5, linestyle="--", alpha=0.5)
    # Trendline
    z = np.polyfit(embds, deltas, 1)
    xs = np.linspace(min(embds), max(embds), 50)
    ax.plot(xs, z[0] * xs + z[1], color=info["color"], alpha=0.4, linewidth=2)
    ax.set_xlabel("emb_d (distance from baseline embedding)")
    ax.set_ylabel("delta (rarity × coherence - baseline)")
    ax.set_title(f"{model}\nPearson r = {r:+.3f}")
    ax.grid(True, alpha=0.2, zorder=1)

plt.suptitle(
    "Delta vs emb_d anti-correlation across the productionized metacog recipe set\n"
    "Heavy register-imposition recipes drive emb_d high but tank delta",
    fontsize=12, y=1.02,
)
plt.tight_layout()
plt.savefig("/home/gabe/metacog/docs/figures/delta-vs-embd-anti-correlation.png", dpi=140, bbox_inches="tight")
print("Saved docs/figures/delta-vs-embd-anti-correlation.png")
