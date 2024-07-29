export default function KaraokeCreate() {
  return (
    <>
      <h1 class="text-6xl font-bold mt-16 mb-8">Create Karaoke</h1>

      <form class="flex flex-col">
        <label>
          <div class="label">
            <span class="label-text">Title</span>
          </div>
          <input
            type="text"
            placeholder="Zankoku na Tenshi no These"
            class="input input-bordered w-full max-w-xs"
          />
        </label>
      </form>
    </>
  );
}
