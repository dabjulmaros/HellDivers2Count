let chart, data;
//Chart options constant defined in line 56

document.addEventListener("htmx:afterRequest", function (evt) {
  const canvasEle = document.querySelector("canvas#playerChart");
  const playerCounts = document.querySelector("div.flex");
  if (evt.target == canvasEle) {
    data = JSON.parse(evt.detail.target.innerHTML);
    evt.target.innerHTML = "";
    if (chart == undefined) {
      createChart();
      updateChart();
    } else {
      updateChart(false);
    }
    canvasEle.style.width = "90vw";
  } else if (evt.target == playerCounts) {
    const updatedEle = document.querySelector("#updated");
    const updatedDate = new Date(updatedEle.innerText);
    updatedEle.innerText = new Intl.DateTimeFormat("default", {
      hour: "2-digit",
      minute: "2-digit",
    }).format(updatedDate);
  }
});

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
