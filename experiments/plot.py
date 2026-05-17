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


def aggregate_recipes(results_path: Path | None = None) -> dict[str, dict]:
    """Return {recipe_name: {delta, emb_d, n}} averaged across tasks."""
    rows = load_rows(results_path or RESULTS)
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


def plot_cross_model(sonnet: dict[str, dict], codex: dict[str, dict]) -> None:
    """Show Sonnet vs codex on (delta, emb_d) for matched recipes.

    Highlights the cross-model winner (envoy-extreme) and the cross-model
    catastrophe (counterpoint-biblical-duo). Connects matched (recipe x
    model) points with arrows to make the transfer (or anti-transfer) visible.
    """
    fig, ax = plt.subplots(figsize=(11, 8))

    # Recipes present on both backends.
    matched = sorted(set(sonnet) & set(codex))
    pairs = []
    for name in matched:
        s = sonnet[name]
        c = codex[name]
        if s.get("control") or c.get("control"):
            continue
        pairs.append((name, s, c))

    # Highlights.
    HIGHLIGHTS = {
        "envoy-extreme": ("#1f77b4", "cross-model winner"),
        "envoy-extreme-alt2": ("#1f77b4", None),
        "counterpoint-biblical-duo": ("#d62728", "cross-model catastrophe"),
        "chorus-plus-disjunction": ("#9467bd", "Sonnet champion"),
    }

    for name, s, c in pairs:
        color, _ = HIGHLIGHTS.get(name, ("#999999", None))
        # Connect Sonnet (circle) -> codex (triangle) with an arrow.
        ax.annotate(
            "", xy=(c["delta"], c["emb_dist"]),
            xytext=(s["delta"], s["emb_dist"]),
            arrowprops=dict(arrowstyle="->", color=color, alpha=0.55, lw=1.4),
            zorder=2,
        )
        # Sonnet point.
        ax.scatter([s["delta"]], [s["emb_dist"]], c=color, s=110,
                   marker="o", edgecolors="black", linewidths=0.9, zorder=4)
        # Codex point.
        ax.scatter([c["delta"]], [c["emb_dist"]], c=color, s=110,
                   marker="^", edgecolors="black", linewidths=0.9, zorder=4)
        # Label at the midpoint.
        if name in HIGHLIGHTS:
            mx = (s["delta"] + c["delta"]) / 2
            my = (s["emb_dist"] + c["emb_dist"]) / 2
            ax.annotate(name, (mx, my),
                        xytext=(8, 8), textcoords="offset points",
                        fontsize=9, fontweight="bold", color=color)

    # Reference lines.
    ax.axhline(0, color="#bbb", linestyle=":", linewidth=1, zorder=0)
    ax.axvline(0, color="#bbb", linestyle=":", linewidth=1, zorder=0)

    ax.set_xlabel("delta (rarity-weighted citation density above NULL baseline)",
                  fontsize=11)
    ax.set_ylabel("emb_d (cosine distance from per-task NULL embedding centroid)",
                  fontsize=11)
    ax.set_title(
        "Cross-model: Sonnet (circle) -> codex (triangle), matched recipes\n"
        "Author-becomes transfer (envoy-extreme); register-shifts and disjunction don't",
        fontsize=12,
    )

    legend_handles = [
        mpatches.Patch(color="#1f77b4", label="cross-model winner (envoy-extreme variants)"),
        mpatches.Patch(color="#d62728", label="cross-model catastrophe (CBD)"),
        mpatches.Patch(color="#9467bd", label="Sonnet-specific champion (chorus-plus-disjunction)"),
        mpatches.Patch(color="#999999", label="other matched recipes"),
    ]
    ax.legend(handles=legend_handles, loc="upper left", fontsize=9)
    ax.grid(True, alpha=0.3)
    fig.tight_layout()
    out = FIGURES_DIR / "cross-model-pareto.png"
    fig.savefig(out, dpi=140)
    print(f"wrote {out}")
    plt.close(fig)


def plot_pareto_current(sonnet: dict[str, dict], codex: dict[str, dict]) -> None:
    """Updated Pareto frontier with post-v6.5.1 winners: anchor-duo,
    duo-disjunction, commitment-excerpt-biblical, and the cross-model
    anchor-duo champion on codex."""
    fig, ax = plt.subplots(figsize=(12, 8.5))

    # All Sonnet recipes as small grey dots (with N>=10).
    others_d, others_e = [], []
    HIGHLIGHTS_S = {
        "chorus-plus-disjunction": "antinomy",
        "trinity-prepended-register": "envoy",
        "envoy-extreme": "envoy-extreme",
        "counterpoint-duo": "counterpoint",
        "envoy-biblical": "envoy-biblical",
        "B1-commitment-disjunction-duo": "duo-disjunction",
        "R9B2-commitment-double-excerpt": "anchor-duo (Sonnet)",
        "commitment-excerpt-biblical": "commit-excerpt-bib",
        "commitment-trio-biblical": "commit-trio-bib",
        "excerpt-biblical-trio": "excerpt-trio",
        "excerpt-biblical-duo": "excerpt-bib-duo",
        "silence-double-excerpt": "silence-double",
        "manifold-cascade": "manifold-cascade",
    }
    for name, r in sonnet.items():
        if r["n"] < 10 or r["control"]:
            continue
        if name in HIGHLIGHTS_S:
            continue
        others_d.append(r["delta"])
        others_e.append(r["emb_dist"])
    ax.scatter(others_d, others_e, c="#cccccc", s=18, alpha=0.5,
               label="other Sonnet recipes", zorder=1)

    # Sonnet productionized stratagems as larger blue dots, labeled.
    PRODUCTIONIZED_S = {
        "chorus-plus-disjunction": "antinomy",
        "trinity-prepended-register": "envoy",
        "envoy-extreme": "envoy-extreme",
        "counterpoint-duo": "counterpoint",
        "B1-commitment-disjunction-duo": "duo-disjunction",
        "R9B2-commitment-double-excerpt": "anchor-duo",
    }
    for recipe, label in PRODUCTIONIZED_S.items():
        if recipe not in sonnet:
            continue
        r = sonnet[recipe]
        ax.scatter([r["delta"]], [r["emb_dist"]], c="#1f77b4", s=180,
                   edgecolors="black", linewidths=1.4, zorder=4)
        # Place labels carefully for the densely-packed frontier.
        offset = {"antinomy": (8, 4), "envoy": (8, 4),
                  "envoy-extreme": (-90, 4), "counterpoint": (8, -12),
                  "duo-disjunction": (8, 4), "anchor-duo": (8, -14)}
        dx, dy = offset.get(label, (8, 6))
        ax.annotate(label, (r["delta"], r["emb_dist"]),
                    xytext=(dx, dy), textcoords="offset points",
                    fontsize=11, fontweight="bold", color="#1f77b4")

    # Additional Pareto-interesting Sonnet recipes (non-productionized winners).
    EXTRAS_S = {
        "commitment-excerpt-biblical": ("emb_d champion (Sonnet)", "#2ca02c", (8, 4)),
        "envoy-biblical": ("envoy-biblical", "#2ca02c", (-90, -12)),
    }
    for recipe, (label, color, off) in EXTRAS_S.items():
        if recipe not in sonnet:
            continue
        r = sonnet[recipe]
        ax.scatter([r["delta"]], [r["emb_dist"]], c=color, s=160,
                   edgecolors="black", linewidths=1.4, zorder=4, marker="D")
        ax.annotate(label, (r["delta"], r["emb_dist"]),
                    xytext=off, textcoords="offset points",
                    fontsize=10, fontweight="bold", color=color)

    # Codex anchor-duo and envoy-extreme as separate red markers.
    CODEX_HIGHLIGHTS = {
        "R9B2-commitment-double-excerpt": "anchor-duo (codex)",
        "envoy-extreme": "envoy-extreme (codex)",
    }
    for recipe, label in CODEX_HIGHLIGHTS.items():
        if recipe not in codex:
            continue
        r = codex[recipe]
        ax.scatter([r["delta"]], [r["emb_dist"]], c="#d62728", s=180,
                   edgecolors="black", linewidths=1.4, zorder=4, marker="s")
        offset = {"anchor-duo (codex)": (8, -14), "envoy-extreme (codex)": (-130, 4)}
        dx, dy = offset.get(label, (8, 6))
        ax.annotate(label, (r["delta"], r["emb_dist"]),
                    xytext=(dx, dy), textcoords="offset points",
                    fontsize=11, fontweight="bold", color="#d62728")

    # Draw the Sonnet Pareto frontier as a connecting line through champions.
    pareto_recipes_s = [
        "commitment-excerpt-biblical",
        "B1-commitment-disjunction-duo",
        "R9B2-commitment-double-excerpt",
        "chorus-plus-disjunction",
    ]
    pareto_points = sorted(
        [(sonnet[r]["delta"], sonnet[r]["emb_dist"])
         for r in pareto_recipes_s if r in sonnet],
        key=lambda p: p[0])
    if pareto_points:
        ax.plot([p[0] for p in pareto_points],
                [p[1] for p in pareto_points],
                color="#1f77b4", linestyle="--", linewidth=1.5,
                alpha=0.4, zorder=2, label="Sonnet Pareto frontier")

    # NULL noise floor reference line.
    ax.axhline(0.090, color="#999", linestyle=":", linewidth=1, zorder=0)
    ax.text(0.42, 0.092, "NULL noise floor (emb_d ~0.09)",
            fontsize=8, color="#666", va="bottom", ha="right")
    ax.axvline(0, color="#999", linestyle=":", linewidth=1, zorder=0)

    ax.set_xlabel("delta  (rarity-weighted citation density above NULL baseline)",
                  fontsize=11)
    ax.set_ylabel("emb_d  (cosine distance from NULL embedding centroid)",
                  fontsize=11)
    ax.set_title("Current Pareto frontier: 8 productionized stratagems\n"
                 "(v6.0.0 -> v6.7.2, 12 rounds, ~80 recipes, ~3500 trials; "
                 "Sonnet blue, Codex red square)",
                 fontsize=12)

    handles = [
        mpatches.Patch(color="#1f77b4", label="productionized stratagem (Sonnet)"),
        mpatches.Patch(color="#d62728", label="cross-model probe (Codex)"),
        mpatches.Patch(color="#2ca02c", label="emb_d champion / register variant"),
        mpatches.Patch(color="#cccccc", label="other Sonnet recipe"),
    ]
    ax.legend(handles=handles, loc="lower right", fontsize=10)
    ax.grid(True, alpha=0.3)
    fig.tight_layout()
    out = FIGURES_DIR / "pareto-frontier-current.png"
    fig.savefig(out, dpi=140)
    print(f"wrote {out}")
    plt.close(fig)


def plot_pareto_three_models(
    sonnet: dict[str, dict],
    opus: dict[str, dict],
    codex: dict[str, dict],
) -> None:
    """Three-panel Pareto frontier: Sonnet / Opus / Codex.

    Highlights v6.8.0 productionized stratagems and the model-specific
    architectural fingerprints (Sonnet rewards scaffolding, Opus rewards
    bare anchors, codex needs occult-anchors + author-extremity).
    """
    fig, axes = plt.subplots(1, 3, figsize=(18, 6), sharey=False)

    # Recipe name -> display label. Used across all three panels;
    # not every recipe is present in every model.
    LABELS = {
        "null": "NULL",
        "chorus": "chorus",
        "chorus-plus-disjunction": "antinomy",
        "envoy-extreme": "envoy-extreme",
        "R9B2-commitment-double-excerpt": "anchor-duo (R9B2)",
        "R12-anchor-duo-occult": "anchor-duo",
        "R12-sigil-name-commitment": "sigil",
        "R12-grimoire-register": "grimoire",
        "R12-occult-cosmologists": "occult-extreme",
        "R14-anchor-duo-name": "anchor-duo-name",
        "R19-chord-anchor": "chord-anchor",
        "R19-witness-anchor": "witness-anchor",
        "R20-chord-witness-anchor": "chord-witness-anchor",
        "R21-chord-of-becomes": "chord-of-becomes",
        "R22-chord-anchor-bare": "chord-anchor-bare",
        "R23-chord-extreme-occult": "chord-extreme",
    }

    # Per-model recipes worth highlighting and their offsets for label placement.
    # Format: name -> (color, offset_xy)
    PANEL_HIGHLIGHTS = {
        "Sonnet (claude-sonnet-4-6)": {
            "data": sonnet,
            "highlights": {
                "R14-anchor-duo-name": ("#d62728", (8, 4)),     # delta champion
                "R19-chord-anchor": ("#1f77b4", (8, 4)),        # Pareto chord
                "R22-chord-anchor-bare": ("#1f77b4", (8, -14)),
                "R12-sigil-name-commitment": ("#9467bd", (-90, 4)),
                "R12-anchor-duo-occult": ("#999999", (8, -14)),
                "chorus-plus-disjunction": ("#999999", (8, 4)),
                "envoy-extreme": ("#2ca02c", (-90, -12)),
            },
            "headline": "delta ceiling: anchor-duo-name +0.359 (N=30)",
        },
        "Opus (claude-opus-4-7)": {
            "data": opus,
            "highlights": {
                "R12-anchor-duo-occult": ("#d62728", (8, -14)),    # delta champion
                "R19-chord-anchor": ("#1f77b4", (8, 4)),           # dual-axis
                "R22-chord-anchor-bare": ("#1f77b4", (-105, -14)),
                "R20-chord-witness-anchor": ("#9467bd", (8, 4)),   # emb_d champion
                "R12-sigil-name-commitment": ("#999999", (8, -14)),
                "R14-anchor-duo-name": ("#999999", (8, 4)),
                "envoy-extreme": ("#2ca02c", (-100, 4)),
                "chorus-plus-disjunction": ("#999999", (8, 4)),
            },
            "headline": "delta ceiling: anchor-duo +0.596 (N=20); emb_d ceiling 0.350",
        },
        "Codex (gpt-5.5)": {
            "data": codex,
            "highlights": {
                "R12-occult-cosmologists": ("#d62728", (8, 4)),   # only stable winner
                "envoy-extreme": ("#2ca02c", (8, -14)),
                "R12-anchor-duo-occult": ("#999999", (8, 4)),     # tank
                "R19-chord-anchor": ("#999999", (8, -14)),        # tank
                "R12-sigil-name-commitment": ("#999999", (-100, 4)),
            },
            "headline": "envoy-extreme transfers (+0.283 N=36); chord recipes fail",
        },
    }

    for ax, (title, panel) in zip(axes, PANEL_HIGHLIGHTS.items()):
        data = panel["data"]
        highlights = panel["highlights"]

        # Plot every recipe with N>=10 as a grey background dot.
        bg_d, bg_e = [], []
        for name, r in data.items():
            if r["n"] < 10 or r["control"]:
                continue
            if name in highlights:
                continue
            bg_d.append(r["delta"])
            bg_e.append(r["emb_dist"])
        ax.scatter(bg_d, bg_e, c="#dddddd", s=14, alpha=0.6, zorder=1)

        # NULL baseline at origin.
        if "null" in data:
            n = data["null"]
            ax.scatter([n["delta"]], [n["emb_dist"]], c="black", s=80,
                       marker="x", zorder=3)
            ax.annotate("NULL", (n["delta"], n["emb_dist"]),
                        xytext=(6, 6), textcoords="offset points",
                        fontsize=9, color="black")

        # Highlighted recipes.
        for recipe, (color, offset) in highlights.items():
            if recipe not in data:
                continue
            r = data[recipe]
            ax.scatter([r["delta"]], [r["emb_dist"]], c=color, s=130,
                       edgecolors="black", linewidths=1.0, zorder=4)
            label = LABELS.get(recipe, recipe)
            ax.annotate(label, (r["delta"], r["emb_dist"]),
                        xytext=offset, textcoords="offset points",
                        fontsize=9, fontweight="bold", color=color)

        ax.axhline(0, color="#bbb", linestyle=":", linewidth=1, zorder=0)
        ax.axvline(0, color="#bbb", linestyle=":", linewidth=1, zorder=0)
        ax.set_xlabel("delta (rarity x coherence above NULL)", fontsize=10)
        if ax is axes[0]:
            ax.set_ylabel("emb_d (cosine distance from NULL centroid)",
                          fontsize=10)
        ax.set_title(f"{title}\n{panel['headline']}", fontsize=11)
        ax.grid(True, alpha=0.25)

    # Shared legend (drawn outside the rightmost panel).
    legend_handles = [
        mpatches.Patch(color="#d62728", label="delta champion"),
        mpatches.Patch(color="#1f77b4", label="chord-anchor / Pareto-emb_d"),
        mpatches.Patch(color="#9467bd", label="other notable"),
        mpatches.Patch(color="#2ca02c", label="cross-model transfer"),
        mpatches.Patch(color="#dddddd", label="other recipe N>=10"),
    ]
    fig.legend(handles=legend_handles, loc="lower center", ncol=5,
               fontsize=9, bbox_to_anchor=(0.5, -0.02))

    fig.suptitle(
        "Pareto frontier across three generators "
        "(18 primitives, 27 stratagems, v6.8.1, all recipes shown at N>=20)",
        fontsize=13, y=1.00,
    )
    fig.tight_layout(rect=(0, 0.03, 1, 0.97))
    out = FIGURES_DIR / "pareto-three-models.png"
    fig.savefig(out, dpi=140, bbox_inches="tight")
    print(f"wrote {out}")
    plt.close(fig)


def main():
    recipes = aggregate_recipes()
    print(f"loaded {len(recipes)} recipes")
    plot_pareto(recipes)
    plot_matrix(recipes)
    plot_register_triangulation(recipes)
    plot_arc(recipes)

    codex_path = Path(__file__).parent / "codex_results.tsv"
    opus_path = Path(__file__).parent / "opus_results.tsv"
    codex = aggregate_recipes(codex_path) if codex_path.exists() else {}
    opus = aggregate_recipes(opus_path) if opus_path.exists() else {}
    if codex:
        print(f"loaded {len(codex)} codex recipes")
        plot_cross_model(recipes, codex)
        plot_pareto_current(recipes, codex)
    if opus:
        print(f"loaded {len(opus)} opus recipes")
    if codex and opus:
        plot_pareto_three_models(recipes, opus, codex)


if __name__ == "__main__":
    main()
