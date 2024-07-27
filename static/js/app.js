// htmx.logAll();
const draggables = document.querySelectorAll(".task");
const droppables = document.querySelectorAll(".swim-lane");

document.body.addEventListener('dragend', function(evt){
    console.log(evt.toElement.id);
    console.log("fire ajax here with status=" + evt.toElement.parentElement.id);
    var body = {
        uuid: evt.toElement.id,
        status: evt.toElement.parentElement.id.replace("-lane", "")
    };
    htmx.ajax('POST', '/todos/'+ evt.toElement.id, {values: body, target: "#myDiv"});
});

document.body.addEventListener('htmx:load', function(evt){
    console.log("body on load");
    updateDraggbles(draggables);
    updateDroppables(droppables);

});

document.body.addEventListener('htmx:afterSwap', function(evt) {
    console.log(evt)
    if(evt.detail.xhr.status === 200){
        updateDraggbles(document.querySelectorAll(".task"));
        console.log(evt.detail.target)
    }
});

function updateDraggbles(draggables){
    draggables.forEach((task) => {
        task.addEventListener("dragstart", () => {
          task.classList.add("is-dragging");
        });
        task.addEventListener("dragend", () => {
          task.classList.remove("is-dragging");
        });
    });
}

function updateDroppables(droppables){
    droppables.forEach((zone) => {
        zone.addEventListener("dragover", (e) => {
          e.preventDefault();
          const bottomTask = insertAboveTask(zone, e.clientY);
          const curTask = document.querySelector(".is-dragging");

          if (!bottomTask) {
            zone.appendChild(curTask);
          } else {
            zone.insertBefore(curTask, bottomTask);
          }
        });
    });
}

const insertAboveTask = (zone, mouseY) => {
    const els = zone.querySelectorAll(".task:not(.is-dragging)");

    let closestTask = null;
    let closestOffset = Number.NEGATIVE_INFINITY;
    els.forEach((task) => {
      const { top } = task.getBoundingClientRect();

      const offset = mouseY - top;

      if (offset < 0 && offset > closestOffset) {
        closestOffset = offset;
        closestTask = task;
      }
    });
    return closestTask;
  };