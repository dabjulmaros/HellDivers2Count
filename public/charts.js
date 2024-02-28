document.addEventListener("htmx:afterRequest", function (evt) {
  if (evt.target == document.querySelector("canvas#playerChart")) {
    data = JSON.parse(evt.detail.target.innerHTML);
    evt.target.innerHTML = "";
    if (chart == undefined) {
      createChart();
      updateChart();
    } else {
      updateChart(false);
    }
  }
});

const chartOptions = {
  elements: {
    point: {
      pointStyle: true,
    },
  },
  scales: {
    y: {
      title: {
        color: "#171d25",
      },
      beginAtZero: true,
      grid: {
        color: "#171d25",
        tickColor: "#171d25",
      },
    },
    x: {
      title: {
        color: "#171d25",
      },
      grid: {
        color: "#171d25",
        tickColor: "#171d25",
      },
      ticks: {
        display: false,
      },
    },
  },
  plugins: {
    legend: {
      display: false,
    },
  },
};
let chart, data;
function createChart() {
  const ctx = document.getElementById("playerChart");
  chart = new Chart(ctx, {
    type: "line",
    data: {
      labels: [],
      datasets: [],
    },
    options: chartOptions,
  });
}

function updateChart(animate = true) {
  chart.data.labels = [
    ...data.map((e) => new Date(e["Updated"]).toLocaleString()),
  ];
  chart.data.datasets = [
    {
      label: "# of Steam Players",
      data: data.map((e) => e["Count"]),
      borderColor: "#ffe80a",
      backgroundColor: "#ffe80a40",
      fill: true,
      tension: 0.5,
    },
  ];
  chart.update(animate ? "" : "none");
}
