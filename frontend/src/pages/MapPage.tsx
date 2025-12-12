import React, { useEffect, useRef, useState } from 'react';
import mapboxgl from 'mapbox-gl';
import 'mapbox-gl/dist/mapbox-gl.css';
import { useQuery } from '@tanstack/react-query';
import { api } from '@/services/api';
import { NODE_COLORS, NODE_ICONS, type NodeType } from '@/types';
import { Card, Badge, Button } from '@/components/ui';
import { Layers, ZoomIn, ZoomOut, Locate, Maximize2 } from 'lucide-react';
import { cn } from '@/lib/utils';

// Set Mapbox access token from environment
mapboxgl.accessToken = import.meta.env.VITE_MAPBOX_TOKEN || '';

export function MapPage() {
    const mapContainer = useRef<HTMLDivElement>(null);
    const map = useRef<mapboxgl.Map | null>(null);
    const [mapLoaded, setMapLoaded] = useState(false);
    const [selectedType, setSelectedType] = useState<NodeType | 'ALL'>('ALL');

    // Fetch GeoJSON data
    const { data: nodesGeoJSON } = useQuery({
        queryKey: ['nodes-geojson'],
        queryFn: () => api.getNodesGeoJSON(),
    });

    // Initialize map
    useEffect(() => {
        if (!mapContainer.current || map.current) return;

        map.current = new mapboxgl.Map({
            container: mapContainer.current,
            style: 'mapbox://styles/mapbox/dark-v11',
            center: [106.8456, -6.2088], // Jakarta center
            zoom: 12,
            pitch: 0,
        });

        map.current.on('load', () => {
            setMapLoaded(true);
        });

        // Add navigation controls
        map.current.addControl(new mapboxgl.NavigationControl(), 'bottom-right');

        return () => {
            if (map.current) {
                map.current.remove();
                map.current = null;
            }
        };
    }, []);

    // Add nodes layer when data is ready
    useEffect(() => {
        if (!map.current || !mapLoaded || !nodesGeoJSON) return;

        // Remove existing source/layer if present
        if (map.current.getLayer('nodes-layer')) {
            map.current.removeLayer('nodes-layer');
        }
        if (map.current.getLayer('nodes-labels')) {
            map.current.removeLayer('nodes-labels');
        }
        if (map.current.getSource('nodes')) {
            map.current.removeSource('nodes');
        }

        // Add GeoJSON source
        map.current.addSource('nodes', {
            type: 'geojson',
            data: nodesGeoJSON as GeoJSON.FeatureCollection,
        });

        // Add circle layer for nodes
        map.current.addLayer({
            id: 'nodes-layer',
            type: 'circle',
            source: 'nodes',
            paint: {
                'circle-radius': [
                    'match',
                    ['get', 'type'],
                    'OLT', 16,
                    'ODC', 12,
                    'ODP', 8,
                    'CLOSURE', 8,
                    'POLE', 6,
                    'CUSTOMER', 6,
                    8,
                ],
                'circle-color': [
                    'match',
                    ['get', 'type'],
                    'OLT', NODE_COLORS.OLT,
                    'ODC', NODE_COLORS.ODC,
                    'ODP', NODE_COLORS.ODP,
                    'CLOSURE', NODE_COLORS.CLOSURE,
                    'POLE', NODE_COLORS.POLE,
                    'CUSTOMER', NODE_COLORS.CUSTOMER,
                    '#666666',
                ],
                'circle-stroke-width': 2,
                'circle-stroke-color': '#ffffff',
                'circle-opacity': 0.9,
            },
            filter: selectedType === 'ALL' ? ['has', 'type'] : ['==', ['get', 'type'], selectedType],
        });

        // Add labels for OLT and ODC
        map.current.addLayer({
            id: 'nodes-labels',
            type: 'symbol',
            source: 'nodes',
            layout: {
                'text-field': ['get', 'name'],
                'text-size': 11,
                'text-offset': [0, 1.5],
                'text-anchor': 'top',
            },
            paint: {
                'text-color': '#ffffff',
                'text-halo-color': '#000000',
                'text-halo-width': 1,
            },
            filter: ['in', ['get', 'type'], ['literal', ['OLT', 'ODC']]],
        });

        // Add click handler for popups
        map.current.on('click', 'nodes-layer', (e) => {
            if (!e.features || e.features.length === 0) return;

            const feature = e.features[0];
            const coordinates = (feature.geometry as GeoJSON.Point).coordinates.slice() as [number, number];
            const props = feature.properties;

            new mapboxgl.Popup({ offset: 15 })
                .setLngLat(coordinates)
                .setHTML(`
          <div class="p-2">
            <h3 class="font-bold text-white text-lg mb-2">${props?.name || 'Unknown'}</h3>
            <div class="space-y-1 text-sm">
              <div class="flex justify-between">
                <span class="text-slate-400">Type:</span>
                <span class="text-white">${props?.type}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-slate-400">Status:</span>
                <span class="text-white">${props?.status}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-slate-400">Capacity:</span>
                <span class="text-white">${props?.used_ports}/${props?.capacity_ports} ports</span>
              </div>
            </div>
          </div>
        `)
                .addTo(map.current!);
        });

        // Change cursor on hover
        map.current.on('mouseenter', 'nodes-layer', () => {
            if (map.current) map.current.getCanvas().style.cursor = 'pointer';
        });

        map.current.on('mouseleave', 'nodes-layer', () => {
            if (map.current) map.current.getCanvas().style.cursor = '';
        });
    }, [mapLoaded, nodesGeoJSON, selectedType]);

    // Filter function
    const handleFilterChange = (type: NodeType | 'ALL') => {
        setSelectedType(type);
        if (map.current && mapLoaded) {
            map.current.setFilter(
                'nodes-layer',
                type === 'ALL' ? ['has', 'type'] : ['==', ['get', 'type'], type]
            );
        }
    };

    // Zoom controls
    const handleZoomIn = () => map.current?.zoomIn();
    const handleZoomOut = () => map.current?.zoomOut();
    const handleResetView = () => {
        map.current?.flyTo({ center: [106.8456, -6.2088], zoom: 12 });
    };

    const nodeTypes: (NodeType | 'ALL')[] = ['ALL', 'OLT', 'ODC', 'ODP', 'CLOSURE', 'POLE', 'CUSTOMER'];

    return (
        <div className="relative h-screen">
            {/* Map Container */}
            <div ref={mapContainer} className="absolute inset-0" />

            {/* Map Controls */}
            <div className="absolute top-4 left-4 z-10">
                <Card variant="glass" className="p-4">
                    <div className="flex items-center gap-2 mb-3">
                        <Layers className="h-5 w-5 text-sky-400" />
                        <h3 className="font-semibold text-white">Filter Nodes</h3>
                    </div>
                    <div className="flex flex-wrap gap-2">
                        {nodeTypes.map((type) => (
                            <button
                                key={type}
                                onClick={() => handleFilterChange(type)}
                                className={cn(
                                    'px-3 py-1.5 rounded-lg text-sm font-medium transition-all',
                                    selectedType === type
                                        ? 'bg-sky-500 text-white'
                                        : 'bg-slate-700/50 text-slate-300 hover:bg-slate-600/50'
                                )}
                            >
                                {type === 'ALL' ? 'All' : type}
                                {type !== 'ALL' && (
                                    <span className="ml-1.5">{NODE_ICONS[type]}</span>
                                )}
                            </button>
                        ))}
                    </div>
                </Card>
            </div>

            {/* Zoom Controls */}
            <div className="absolute bottom-24 right-4 z-10 flex flex-col gap-2">
                <Button size="icon" variant="outline" onClick={handleZoomIn}>
                    <ZoomIn className="h-5 w-5" />
                </Button>
                <Button size="icon" variant="outline" onClick={handleZoomOut}>
                    <ZoomOut className="h-5 w-5" />
                </Button>
                <Button size="icon" variant="outline" onClick={handleResetView}>
                    <Maximize2 className="h-5 w-5" />
                </Button>
            </div>

            {/* Legend */}
            <div className="absolute bottom-4 left-4 z-10">
                <Card variant="glass" className="p-4">
                    <h4 className="text-sm font-semibold text-white mb-3">Legend</h4>
                    <div className="grid grid-cols-2 gap-2 text-xs">
                        {(Object.keys(NODE_COLORS) as NodeType[]).map((type) => (
                            <div key={type} className="flex items-center gap-2">
                                <div
                                    className="w-3 h-3 rounded-full border border-white"
                                    style={{ backgroundColor: NODE_COLORS[type] }}
                                />
                                <span className="text-slate-300">{type}</span>
                            </div>
                        ))}
                    </div>
                </Card>
            </div>

            {/* Loading Overlay */}
            {!mapLoaded && (
                <div className="absolute inset-0 flex items-center justify-center bg-slate-900/80 z-20">
                    <div className="text-center">
                        <div className="w-12 h-12 border-4 border-sky-500 border-t-transparent rounded-full animate-spin mx-auto mb-4" />
                        <p className="text-slate-400">Loading map...</p>
                    </div>
                </div>
            )}

            {/* Token Warning */}
            {!import.meta.env.VITE_MAPBOX_TOKEN && (
                <div className="absolute top-4 right-4 z-10">
                    <Card variant="glass" className="p-4 border-amber-500/50">
                        <Badge variant="warning" className="mb-2">Mapbox Token Required</Badge>
                        <p className="text-sm text-slate-400">
                            Add your Mapbox token to <code className="text-sky-400">.env</code>
                        </p>
                    </Card>
                </div>
            )}
        </div>
    );
}
