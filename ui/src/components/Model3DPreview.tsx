import { Canvas, useThree } from "@react-three/fiber";
import { PerspectiveCamera, OrbitControls as DreiOrbitControls } from "@react-three/drei";
import { useEffect, useState, useRef } from "react";
import { STLLoader, ThreeMFLoader } from "three-stdlib";
import * as THREE from "three";

type Model3DPreviewProps = {
  modelData: string; // base64 string
  color?: string;
  width?: number;
  height?: number;
  fileType?: "stl" | "3mf";
};

const Controls = () => {
  const { invalidate } = useThree();
  return (
    <DreiOrbitControls
      enableDamping
      dampingFactor={0.1}
      rotateSpeed={0.5}
      target={[0, 0, 0]}
      onChange={() => invalidate()}
    />
  );
};

const ModelViewer = ({
  modelData,
  color = "",
  fileType = "stl",
}: {
  modelData: string;
  color: string;
  fileType?: "stl" | "3mf";
}) => {
  const [parsed, setParsed] = useState<{
    object3d: THREE.Object3D | null;
    boundingBox: THREE.Box3 | null;
  } | null>(null);

  const meshRef = useRef<THREE.Object3D>(null);
  const cameraRef = useRef<THREE.PerspectiveCamera>(null);
  const { size, invalidate } = useThree();

  useEffect(() => {
    const parseModel = () => {
      try {
        const binary = atob(modelData);
        const buffer = new Uint8Array(binary.length);
        for (let i = 0; i < binary.length; i++) buffer[i] = binary.charCodeAt(i);

        if (fileType === "3mf") {
          const loader = new ThreeMFLoader();
          const group = loader.parse(buffer.buffer);
          const box = new THREE.Box3().setFromObject(group);
          const center = box.getCenter(new THREE.Vector3());
          group.position.sub(center);
          setParsed({
            object3d: group,
            boundingBox: box.translate(center.negate()),
          });
        } else {
          const loader = new STLLoader();
          const geometry = loader.parse(buffer.buffer);
          geometry.computeVertexNormals();
          geometry.computeBoundingBox();
          const center = new THREE.Vector3();
          if (geometry.boundingBox) {
            geometry.boundingBox.getCenter(center);
            geometry.translate(-center.x, -center.y, -center.z);
          }
          const mesh = new THREE.Mesh(
            geometry,
            new THREE.MeshLambertMaterial({ color: color || "#aaaaaa" })
          );
          mesh.castShadow = false;
          mesh.receiveShadow = false;
          setParsed({
            object3d: mesh,
            boundingBox: geometry.boundingBox,
          });
        }
      } catch (error) {
        console.error("Error parsing model", error);
        setParsed({ object3d: null, boundingBox: null });
      }
    };

    parseModel();
  }, [modelData, fileType, color]);

  useEffect(() => {
    if (!cameraRef.current || !parsed?.boundingBox) return;

    const camera = cameraRef.current;
    const sizeVec = new THREE.Vector3();
    parsed.boundingBox.getSize(sizeVec);
    const maxDim = Math.max(sizeVec.x, sizeVec.y, sizeVec.z);
    const fov = (camera.fov * Math.PI) / 180;
    const distance = (maxDim / (2 * Math.tan(fov / 2))) * 1.5;

    camera.position.set(-distance, distance, distance);
    camera.near = 0.01;
    camera.far = distance * 10;
    camera.lookAt(0, 0, 0);
    camera.updateProjectionMatrix();
    invalidate();
  }, [parsed?.boundingBox]);

  if (!parsed || !parsed.object3d) return null;

  return (
    <>
      <PerspectiveCamera
        ref={cameraRef}
        makeDefault
        fov={45}
        aspect={size.width / size.height}
      />
      <ambientLight intensity={0.6} />
      <directionalLight position={[5, 5, 5]} intensity={0.9} />
      <primitive object={parsed.object3d} ref={meshRef} />
      <Controls />
    </>
  );
};

export const Model3DPreview = ({
  modelData,
  color = "",
  width = 400,
  height = 300,
  fileType = "stl",
}: Model3DPreviewProps) => (
  <Canvas
    style={{ width, height }}
    shadows={false}
    frameloop="demand"
    camera={{ near: 0.01, far: 1000 }}
  >
    <ModelViewer modelData={modelData} color={color} fileType={fileType} />
  </Canvas>
);
