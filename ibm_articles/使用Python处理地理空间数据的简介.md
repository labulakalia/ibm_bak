# 使用 Python 处理地理空间数据的简介
先介绍地理空间数据及其类型、矢量和栅格，然后学习如何使用 Python 处理地理空间数据

**标签:** Python,Python Data Analysis Library,分析,数据科学

[原文链接](https://developer.ibm.com/zh/articles/introduction-to-geospatial-data-using-python/)

[Margriet Groenendijk](https://developer.ibm.com/zh/profiles/mgroenen), [Samaya Madhavan](https://developer.ibm.com/zh/profiles/smadhava)

发布: 2020-06-09

* * *

本文已纳入 [使用 Python 执行数据分析](https://developer.ibm.com/zh/series/learning-path-data-analysis-using-python/) 学习路径。

级别主题类型100[在 Python 中使用 pandas 执行数据分析](https://developer.ibm.com/zh/technologies/data-science/tutorials/data-analysis-in-python-using-pandas)教程**101****[使用 Python 处理地理空间数据的简介](https://developer.ibm.com/zh/articles/introduction-to-geospatial-data-using-python/)**文章201a[在 Python 中处理地理空间矢量数据](https://developer.ibm.com/zh/tutorials/working-with-geospatial-vector-data-in-python/)教程201b[在 Python 中处理地理空间栅格数据](https://developer.ibm.com/zh/tutorials/working-with-geospatial-raster-data-in-python/)教程

## 地理空间数据

地理空间数据是指除其他特性外还包含有关记录地理位置信息的任何数据集。例如，某个数据集既包含多个城市及其人口规模相关信息，同时也包含另外两个带有经度和纬度坐标的列，这个数据集就会被视为地理空间数据。地理空间信息有助于推断出许多额外的信息，例如，查找城市之间的距离、按街区计算平均家庭收入以及创建地图等。

通常，以两种方式表示地理空间数据：矢量数据和栅格数据。本文的其余部分将深入探讨这些表示法及其适用应用程序的细节。

### 矢量数据

_矢量数据_ 通过 x 和 y 坐标表示空间元素。矢量数据的最基本形式是 _点_。两个或两个以上的点形成一条 _线_，三条或三条以上的线形成一个 _多边形_。下图显示了每种类型的矢量数据及其数组表示法。

例如，您可以通过一个点（x 和 y 坐标）定义城市的位置，但实际上城市中不同的形状代表更多信息。有些路可以用线表示，其中包含表示道路起点和终点的两个点。另外，还存在许多个多边形，用来表示任何形式的形状，例如建筑物、区域和城市边界。

#### 数据表示

可以采用各种格式来表示矢量数据。最简单的形式就是在表格中另外增加一个或多个列，用来定义地理空间坐标。更正式的编码格式（如 [GeoJSON](https://geojson.org/)）也会派上用场。GeoJSON 是 [JSON](https://www.json.org/json-en.html) 数据格式的扩展，它包含 _几何特性_，可以是 Point、LineString、Polygon、MultiPoint、MultiLineString 或 MultiPolygon。地理空间数据抽象库 ( [GDAL](https://gdal.org/)) 中描述了其他几个可用于表示地理空间数据的库。

此外，还开发了几个 GDAL 兼容的 Python 软件包，以便于在 Python 中处理地理空间数据。点、线和多边形也可以借助 [Shapely](https://shapely.readthedocs.io/) 描述为对象。通过使用这些 Shapely 对象，您可以探究空间关系，例如 _包含、相交、重叠_ 和 _接触_，如下图所示。

您可以通过几种方式使用 Python 处理地理空间数据。例如，您可以使用 [Fiona](https://pypi.org/project/Fiona/) 加载几何数据，然后将其传递给 Shapely 对象。这就是 [GeoPandas](https://geopandas.org/) 使用的内容。GeoPandas 是一个软件包，借助它可以使处理矢量数据的体验与使用 Pandas 处理表格数据的体验相似。

在指定包含 GeoPandas 地理空间数据的文件类型（csv、geojson 或形状文件）后，它会使用 Fiona 从 GDAL 获取正确的格式。在使用以下代码读取数据后，除 Pandas DataFrame 的功能外，您还将拥有一个 GeoDataFrame，它具有所有地理空间特性，例如，查找最接近的点或计算多边形内的面积和点数。

```
geopandas.read_file()

```

Show moreShow more icon

阅读 [GeoPandas 文档](https://geopandas.org/) 即可了解更多信息。

### 栅格数据

栅格数据是另一种类型的地理空间数据。在观察整个区域的空间信息时，会使用栅格数据。它由行和列的矩阵组成，其中包含与每个单元关联的一些信息。城市的卫星图像就是栅格数据的一个示例，该图像由一个矩阵来表示，矩阵的每个单元中都包含天气信息。当您查看测雨雷达以确定是否需要带雨伞时，您查看的就是此类数据。

#### 数据表示

您可以通过几种方式在 Python 中处理栅格数据。最新的一款用户友好型软件包为 [xarray](http://xarray.pydata.org/en/stable/)，它可读取 [netcdf](https://www.unidata.ucar.edu/software/netcdf/docs/faq.html#What-Is-netCDF) 文件。这是一种二进制数据格式，由多个数组、变量名称的元数据、坐标系、栅格大小和数据作者组成。将文件加载为 `DataArray` 后，只用一个命令（参见以下代码）即可创建地图，类似于 Pandas 和 GeoPandas。

```
da.plot()

```

Show moreShow more icon

对于其他文件格式，可尝试 [rasterio](https://rasterio.readthedocs.io/en/latest/intro.html)，它能够读取许多不同种类的栅格数据文件。

## 坐标系

对于这两种类型的地理空间数据，务必要注意坐标系。地图以规则的坐标网格来表示，但地球并不是一个扁平的矩形。从 3 维地球到 2 维地图的数据转换以多种方式完成，这一过程即称为 [地图投影](https://en.wikipedia.org/wiki/Map_projection)。经常使用的 [墨卡托投影](https://en.wikipedia.org/wiki/Mercator_projection) 与 [摩尔魏特投影](https://en.wikipedia.org/wiki/Mollweide_projection) 看起来有些不同。

在跨多个数据集使用 GeoPandas 时，务必要检查所有数据集的投影是否相同。您可以使用以下代码完成此操作。

```
geodf.crs

```

Show moreShow more icon

当使用栅格数据时，必须注意取决于其位置的值，赤道上的值可能大于极点的值，或者正好相反，这取决于所使用的投影。

## 地理空间数据样本

地理空间数据的部分来源如下：

- 来自 [Natural Earth](https://www.naturalearthdata.com/downloads/50m-cultural-vectors/50m-admin-0-countries-2/) 的所有国家或地区的多边形
- 来自 [世界银行](https://datacatalog.worldbank.org/dataset/major-rivers-world) 的全球河流
- 来自 [CRU](https://crudata.uea.ac.uk/cru/data/hrg/) 的全球历史温度数据
- 来自 [NASA](https://firms.modaps.eosdis.nasa.gov/active_fire/#firms-shapefile) 的仍在燃烧的火灾卫星地图数据
- 来自 [NASA](https://modis.gsfc.nasa.gov/data/dataprod/mod12.php) 的土地覆盖数据

在小型项目中探究数据是了解更多信息的最佳方法。尝试通过回答关于这些数据集的问题来探究这些数据集。例如，探究哪个国家或地区的河流最多，或者尝试通过创建这些地区的变化图和时间序列图来确定温度上升最高的地方。

## 结束语

本文为您简要介绍了地理空间数据及其类型、矢量和栅格。同时还简要讨论了如何使用 Python 处理地理空间数据。本系列中的后续教程将更深入地介绍如何使用 Python 处理矢量数据和栅格数据。这里提供了样本地理空间数据集，您可由此开始探究数据并了解更多信息。本文已纳入 [使用 Python 执行数据分析](https://developer.ibm.com/zh/series/learning-path-data-analysis-using-python/) 学习路径。要继续学习，可查看下一步： [在 Python 中处理地理空间矢量数据](https://developer.ibm.com/zh/tutorials/working-with-geospatial-vector-data-in-python/)。

本文翻译自： [Introduction to geospatial data using Python](https://developer.ibm.com/articles/introduction-to-geospatial-data-using-python/)（2020-03-05）