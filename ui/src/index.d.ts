declare module "@jellyfin/libass-wasm" {
  interface OptionsBase {
    /**
     * The video element to attach listeners to
     */
    video?: HTMLVideoElement;

    /**
     * The URL of the worker
     * @default `subtitles-octopus-worker.js`
     */
    workerUrl?: string;

    /**
     * The URL of the legacy worker
     * @default `subtitles-octopus-worker-legacy.js`
     */
    legacyWorkerUrl?: string;

    /**
     * An array of links to the fonts used in the subtitle
     */
    fonts?: string[];
  }

  interface OptionsWithSubUrl extends OptionsBase {
    /** The URL of the subtitle file to play. (Require either subUrl or subContent to be specified) */
    subUrl: string;
  }

  interface OptionsWithSubContent extends OptionsBase {
    /** The content of the subtitle file to play. (Require either subContent or subUrl to be specified) */
    subContent: string;
  }

  export type Options = OptionsWithSubUrl | OptionsWithSubContent;

  class SubtitlesOctopus {
    constructor(options: Options);

    /**
     * Destroy instance
     */
    dispose(): void;
  }

  export default SubtitlesOctopus;
}
