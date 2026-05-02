"""Generate Pareto-frontier and matrix figures for the v6.0 -> v6.5.1 arc.

Reuses analyze.py for data loading and aggregation.
Outputs PNGs to ../docs/figures/.
"""
from __future__ import annotations

from pathlib import Path

import matplotlib.pyplot as plt
import matplotlib.patches as mpatches

from analyze import (
    load_rows,
    baselines_from_rows,
    embedding_centroids_per_task,
    per_recipe,
    RESULTS,
)

FIGURES_DIR = Path(__file__).parent.parent / "docs" / "figures"
FIGURES_DIR.mkdir(parents=True, exist_ok=True)


def aggregate_recipes() -> dict[str, dict]:
    """Return {recipe_name: {delta, emb_d, n}} averaged across tasks."""
    rows = load_rows(RESULTS)
    baselines = baselines_from_rows(rows)
    centroids = embedding_centroids_per_task(rows)
    out = {}
    for r in per_recipe(rows, baselines, centroids):
        if r["delta"] is None or r["emb_dist"] is None:
            continue
        out[r["recipe"]] = {
            "delta": r["delta"],
            "emb_dist": r["emb_dist"],
            "n": r["n"],
            "control": r["control"],
        }
    return out


# Recipe groupings used across plots.
# Productionized empirical stratagems (mapped to their canonical recipe).
PRODUCTIONIZED = {
    "chorus": "trinity-no-synthesis",
    "trinity": "trinity-manifold",
    "antinomy": "chorus-plus-disjunction",
    "envoy": "trinity-prepended-register",
    "counterpoint": "counterpoint-duo",
}

# 2x3 author-matrix recipes.
MATRIX = {
    "antinomy": {
        "CKW": "chorus-plus-disjunction",
        "MRW": "chorus-plus-disjunction-alt",
        "extreme": "antinomy-extreme",
    },
    "envoy": {
        "CKW": "trinity-prepended-register",
        "MRW": "envoy-alt",
        "extreme": "envoy-extreme",
    },
    "counterpoint": {
        "CKW": "register-chorus-disjunction",
        "MRW": "combined-alt",
        "extreme": "counterpoint-extreme",
    },
}

# Register triangulation recipes (envoy with each register).
REGISTERS = {
    "scientific": "envoy-scientific",
    "Victorian": "trinity-prepended-register",
    "biblical": "envoy-biblical",
}

# Failed compositions for the Pareto plot.
FAILED = {
    "chord-not-fork": "chorus-with-chord-not-fork",
    "glossolalia": "chorus-plus-glossolalia",
    "no-ritual": "antinomy-no-ritual",
    "silence-close": "chorus-silence-instead-of-ritual",
}


def plot_pareto(recipes: dict[str, dict]) -> None:
    """Scatter plot of all recipes on (delta, emb_d), with productionized
    stratagems and failures highlighted."""
    fig, ax = plt.subplots(figsize=(11, 8))

    # All recipes as small grey dots with N>=30.
    others_d, others_e = [], []
    for name, r in recipes.items():
        if r["n"] < 30 or r["control"]:
            continue
        if name in PRODUCTIONIZED.values():
            continue
        if name in FAILED.values():
            continue
        others_d.append(r["delta"])
        others_e.append(r["emb_dist"])
    ax.scatter(others_d, others_e, c="#cccccc", s=24, label="other recipes", zorder=1)

    # Productionized stratagems as labeled larger blue dots.
    for stratagem, recipe in PRODUCTIONIZED.items():
        if recipe not in recipes:
            continue
        r = recipes[recipe]
        ax.scatter([r["delta"]], [r["emb_dist"]], c="#1f77b4", s=140,
                   edgecolors="black", linewidths=1.2, zorder=4)
        ax.annotate(stratagem, (r["delta"], r["emb_dist"]),
                    xytext=(8, 6), textcoords="offset points",
                    fontsize=11, fontweight="bold", color="#1f77b4")

    # Failed compositions as red x.
    for label, recipe in FAILED.items():
        if recipe not in recipes:
            continue
        r = recipes[recipe]
        ax.scatter([r["delta"]], [r["emb_dist"]], c="#d62728", s=80,
                   marker="x", linewidths=2, zorder=3)
        ax.annotate(label, (r["delta"], r["emb_dist"]),
                    xytext=(8, -10), textcoords="offset points",
                    fontsize=9, color="#d62728")

    # envoy-biblical as a special green marker (the structural surprise).
    if "envoy-biblical" in recipes:
        r = recipes["envoy-biblical"]
        ax.scatter([r["delta"]], [r["emb_dist"]], c="#2ca02c", s=140,
                   edgecolors="black", linewidths=1.2, zorder=4)
        ax.annotate("envoy-biblical\n(register variant)", (r["delta"], r["emb_dist"]),
                    xytext=(-90, 6), textcoords="offset points",
                    fontsize=10, fontweight="bold", color="#2ca02c")

    if "envoy-biblical-duo" in recipes:
        r = recipes["envoy-biblical-duo"]
        ax.scatter([r["delta"]], [r["emb_dist"]], c="#2ca02c", s=100,
                   marker="^", edgecolors="black", linewidths=1, zorder=4)
        ax.annotate("biblical-duo", (r["delta"], r["emb_dist"]),
                    xytext=(8, 6), textcoords="offset points",
                    fontsize=9, color="#2ca02c")

    # NULL noise floor reference line.
    ax.axhline(0.090, color="#999", linestyle=":", linewidth=1, zorder=0)
    ax.text(0.42, 0.092, "NULL noise floor (emb_d ~0.09)",
            fontsize=8, color="#666", va="bottom", ha="right")
    ax.axvline(0, color="#999", linestyle=":", linewidth=1, zorder=0)

    ax.set_xlabel("delta  (rarity-weighted citation density above NULL baseline)",
                  fontsize=11)
    ax.set_ylabel("emb_d  (cosine distance from NULL embedding centroid)",
                  fontsize=11)
    ax.set_title("Pareto frontier: where the productionized stratagems sit\n"
                 "(v6.0.0 -> v6.5.1 experimental run, ~50 recipes, ~3000 trials)",
                 fontsize=12)

    handles = [
        mpatches.Patch(color="#1f77b4", label="productionized stratagem"),
        mpatches.Patch(color="#2ca02c", label="biblical-register variant"),
        mpatches.Patch(color="#d62728", label="failed composition"),
        mpatches.Patch(color="#cccccc", label="other recipe"),
    ]
    ax.legend(handles=handles, loc="lower right", fontsize=10)
    ax.grid(True, alpha=0.3)
    fig.tight_layout()
    out = FIGURES_DIR / "pareto-frontier.png"
    fig.savefig(out, dpi=140)
    print(f"wrote {out}")
    plt.close(fig)


def plot_matrix(recipes: dict[str, dict]) -> None:
    """Grouped bar chart of the 2x3 (structure x author) matrix,
    showing both delta and emb_d in side-by-side panels."""
    structures = ["antinomy", "envoy", "counterpoint"]
    authors = ["CKW", "MRW", "extreme"]
    colors = {"CKW": "#1f77b4", "MRW": "#ff7f0e", "extreme": "#9467bd"}

    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(13, 5.5))
    bar_w = 0.26
    x = list(range(len(structures)))

    for i, author in enumerate(authors):
        deltas, embds = [], []
        for s in structures:
            recipe = MATRIX[s][author]
            if recipe in recipes:
                deltas.append(recipes[recipe]["delta"])
                embds.append(recipes[recipe]["emb_dist"])
            else:
                deltas.append(0)
                embds.append(0)
        offset = (i - 1) * bar_w
        ax1.bar([xi + offset for xi in x], deltas, bar_w,
                label=author, color=colors[author], edgecolor="black", linewidth=0.5)
        ax2.bar([xi + offset for xi in x], embds, bar_w,
                label=author, color=colors[author], edgecolor="black", linewidth=0.5)

    for ax, title, ylabel in [
        (ax1, "delta (citation density)", "delta"),
        (ax2, "emb_d (embedding distance)", "emb_d"),
    ]:
        ax.set_xticks(x)
        ax.set_xticklabels(structures, fontsize=11)
        ax.set_ylabel(ylabel, fontsize=11)
        ax.set_title(title, fontsize=12)
        ax.legend(title="author triple", fontsize=9)
        ax.grid(True, alpha=0.3, axis="y")
        ax.axhline(0, color="black", linewidth=0.6)

    fig.suptitle("The 2x3 (structure x author) matrix at N=70+",
                 fontsize=13, y=1.02)
    fig.tight_layout()
    out = FIGURES_DIR / "structure-author-matrix.png"
    fig.savefig(out, dpi=140, bbox_inches="tight")
    print(f"wrote {out}")
    plt.close(fig)


def plot_register_triangulation(recipes: dict[str, dict]) -> None:
    """Bar chart of envoy-{scientific,Victorian,biblical} on both axes."""
    registers = ["scientific", "Victorian", "biblical"]
    colors = ["#1f77b4", "#ff7f0e", "#2ca02c"]

    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(11, 5))
    x = list(range(len(registers)))

    deltas, embds = [], []
    for reg in registers:
        recipe = REGISTERS[reg]
        if recipe in recipes:
            deltas.append(recipes[recipe]["delta"])
            embds.append(recipes[recipe]["emb_dist"])
        else:
            deltas.append(0)
            embds.append(0)

    ax1.bar(x, deltas, color=colors, edgecolor="black", linewidth=0.7)
    ax2.bar(x, embds, color=colors, edgecolor="black", linewidth=0.7)

    for ax, title, ylabel, vals in [
        (ax1, "delta (citation density)", "delta", deltas),
        (ax2, "emb_d (embedding distance)", "emb_d", embds),
    ]:
        ax.set_xticks(x)
        ax.set_xticklabels(registers, fontsize=11)
        ax.set_ylabel(ylabel, fontsize=11)
        ax.set_title(title, fontsize=12)
        ax.grid(True, alpha=0.3, axis="y")
        ax.axhline(0, color="black", linewidth=0.6)
        for xi, v in zip(x, vals):
            ax.text(xi, v, f"{v:+.3f}" if ylabel == "delta" else f"{v:.3f}",
                    ha="center", va="bottom" if v >= 0 else "top", fontsize=10)

    # Annotate biblical-duo as the ceiling-pushed compound on emb_d.
    if "envoy-biblical-duo" in recipes:
        r = recipes["envoy-biblical-duo"]
        ax2.axhline(r["emb_dist"], color="#999", linestyle="--", linewidth=1)
        ax2.text(0.5, r["emb_dist"], f"  envoy-biblical-duo: {r['emb_dist']:.3f}",
                 fontsize=9, color="#666", va="bottom")

    fig.suptitle("Register-target sensitivity (envoy with each register, N>=70)",
                 fontsize=13, y=1.02)
    fig.tight_layout()
    out = FIGURES_DIR / "register-triangulation.png"
    fig.savefig(out, dpi=140, bbox_inches="tight")
    print(f"wrote {out}")
    plt.close(fig)


def plot_arc(recipes: dict[str, dict]) -> None:
    """Show the structural-axis ceiling progression across the v6.0->v6.5.1 arc."""
    milestones = [
        ("v6.1.0\nmanifold-stratagem", "manifold-stratagem", 0),
        ("v6.2.0\nchorus", "trinity-no-synthesis", 1),
        ("v6.4.0\nenvoy-CKW", "trinity-prepended-register", 2),
        ("envoy-extreme", "envoy-extreme", 3),
        ("v6.5.1\ncounterpoint-duo", "counterpoint-duo", 4),
        ("envoy-biblical", "envoy-biblical", 5),
        ("envoy-biblical-duo", "envoy-biblical-duo", 6),
    ]

    fig, ax = plt.subplots(figsize=(12, 6))
    xs, ys, labels = [], [], []
    for label, recipe, idx in milestones:
        if recipe in recipes:
            xs.append(idx)
            ys.append(recipes[recipe]["emb_dist"])
            labels.append(label)

    ax.plot(xs, ys, "o-", color="#1f77b4", linewidth=2, markersize=10,
            markeredgecolor="black", markeredgewidth=1)

    for x, y, label in zip(xs, ys, labels):
        ax.annotate(f"{label}\n{y:.3f}", (x, y),
                    xytext=(0, 12), textcoords="offset points",
                    fontsize=9, ha="center")

    ax.axhline(0.090, color="#999", linestyle=":", linewidth=1)
    ax.text(0.1, 0.092, "NULL noise floor", fontsize=9, color="#666")

    ax.set_xticks(range(len(milestones)))
    ax.set_xticklabels([""] * len(milestones))
    ax.set_xlabel("experimental milestone (chronological)", fontsize=11)
    ax.set_ylabel("emb_d (structural-axis ceiling)", fontsize=11)
    ax.set_title("Structural-axis ceiling progression: v6.0.0 -> v6.5.1\n"
                 "(emb_d of the best-known recipe at each milestone)",
                 fontsize=12)
    ax.grid(True, alpha=0.3, axis="y")
    ax.set_ylim(0.05, 0.35)
    fig.tight_layout()
    out = FIGURES_DIR / "ceiling-progression.png"
    fig.savefig(out, dpi=140)
    print(f"wrote {out}")
    plt.close(fig)


def main():
    recipes = aggregate_recipes()
    print(f"loaded {len(recipes)} recipes")
    plot_pareto(recipes)
    plot_matrix(recipes)
    plot_register_triangulation(recipes)
    plot_arc(recipes)


if __name__ == "__main__":
    main()
