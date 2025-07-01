import { imageState } from '../hooks/useImageProcessor';

export function ProcessedImage() {
  return (
    <div class="mt-4 p-4 border border-gray-700 bg-gray-900 rounded">
      <h3 class="mt-0 text-white text-lg font-semibold">Processed Image Data</h3>
      {imageState.value.processedWidth && (
        <>
          <p class="my-2 text-gray-300">W x H: {imageState.value.processedWidth} px x {imageState.value.processedHeight} px</p>
          <p class="my-2 text-gray-300">Size: {imageState.value.processedSize}</p>
        </>
      )}
      <p>
        <a
          href={imageState.value.downloadUrl ?? undefined}
          target="_blank"
          rel="noopener noreferrer"
          download
          class="inline-block mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
        >
          Download Resized Image
        </a>
      </p>
    </div>
  );
}