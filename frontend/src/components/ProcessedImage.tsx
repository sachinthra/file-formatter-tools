import { imageState } from '../hooks/useImageProcessor';

export function ProcessedImage() {
  return (
    <div class="p-4 border border-gray-700 bg-gray-900 rounded">
      <p>
        <a
          href={imageState.value.downloadUrl ?? undefined}
          target="_blank"
          rel="noopener noreferrer"
          download
          class="inline-block m-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
        >
          Download Resized Image
        </a>
      </p>
    </div>
  );
}